package d3

const dims = 3

const (
	maxEntries = 16
	minEntries = maxEntries * 40 / 100
)

type box struct {
	data     interface{}
	min, max [dims]float64
}

type node struct {
	count int
	boxes [maxEntries + 1]box
}

// BoxTree ...
type BoxTree struct {
	height   int
	root     box
	count    int
	reinsert []box
}

func (r *box) expand(b *box) {
	for i := 0; i < dims; i++ {
		if b.min[i] < r.min[i] {
			r.min[i] = b.min[i]
		}
		if b.max[i] > r.max[i] {
			r.max[i] = b.max[i]
		}
	}
}

func (r *box) area() float64 {
	area := r.max[0] - r.min[0]
	for i := 1; i < dims; i++ {
		area *= r.max[i] - r.min[i]
	}
	return area
}

func (r *box) overlapArea(b *box) float64 {
	area := 1.0
	for i := 0; i < dims; i++ {
		var max, min float64
		if r.max[i] < b.max[i] {
			max = r.max[i]
		} else {
			max = b.max[i]
		}
		if r.min[i] > b.min[i] {
			min = r.min[i]
		} else {
			min = b.min[i]
		}
		if max > min {
			area *= max - min
		} else {
			return 0
		}
	}
	return area
}

func (r *box) enlargedArea(b *box) float64 {
	area := 1.0
	for i := 0; i < len(r.min); i++ {
		if b.max[i] > r.max[i] {
			if b.min[i] < r.min[i] {
				area *= b.max[i] - b.min[i]
			} else {
				area *= b.max[i] - r.min[i]
			}
		} else {
			if b.min[i] < r.min[i] {
				area *= r.max[i] - b.min[i]
			} else {
				area *= r.max[i] - r.min[i]
			}
		}
	}
	return area
}

// Insert inserts an item into the RTree
func (tr *BoxTree) Insert(min, max []float64, value interface{}) {
	var item box
	fit(min, max, value, &item)
	tr.insert(&item)
}

func (tr *BoxTree) insert(item *box) {
	if tr.root.data == nil {
		fit(item.min[:], item.max[:], new(node), &tr.root)
	}
	grown := tr.root.insert(item, tr.height)
	if grown {
		tr.root.expand(item)
	}
	if tr.root.data.(*node).count == maxEntries+1 {
		newRoot := new(node)
		tr.root.splitLargestAxisEdgeSnap(&newRoot.boxes[1])
		newRoot.boxes[0] = tr.root
		newRoot.count = 2
		tr.root.data = newRoot
		tr.root.recalc()
		tr.height++
	}
	tr.count++
}

func (r *box) chooseLeastEnlargement(b *box) int {
	j, jenlargement, jarea := -1, 0.0, 0.0
	n := r.data.(*node)
	for i := 0; i < n.count; i++ {
		var area float64
		if false {
			area = n.boxes[i].area()
		} else {
			// force inline
			area = n.boxes[i].max[0] - n.boxes[i].min[0]
			for j := 1; j < dims; j++ {
				area *= n.boxes[i].max[j] - n.boxes[i].min[j]
			}
		}
		var enlargement float64
		if false {
			enlargement = n.boxes[i].enlargedArea(b) - area
		} else {
			// force inline
			enlargedArea := 1.0
			for j := 0; j < len(n.boxes[i].min); j++ {
				if b.max[j] > n.boxes[i].max[j] {
					if b.min[j] < n.boxes[i].min[j] {
						enlargedArea *= b.max[j] - b.min[j]
					} else {
						enlargedArea *= b.max[j] - n.boxes[i].min[j]
					}
				} else {
					if b.min[j] < n.boxes[i].min[j] {
						enlargedArea *= n.boxes[i].max[j] - b.min[j]
					} else {
						enlargedArea *= n.boxes[i].max[j] - n.boxes[i].min[j]
					}
				}
			}
			enlargement = enlargedArea - area
		}

		if j == -1 || enlargement < jenlargement {
			j, jenlargement, jarea = i, enlargement, area
		} else if enlargement == jenlargement {
			if area < jarea {
				j, jenlargement, jarea = i, enlargement, area
			}
		}
	}
	return j
}

func (r *box) recalc() {
	n := r.data.(*node)
	r.min = n.boxes[0].min
	r.max = n.boxes[0].max
	for i := 1; i < n.count; i++ {
		r.expand(&n.boxes[i])
	}
}

// contains return struct when b is fully contained inside of n
func (r *box) contains(b *box) bool {
	for i := 0; i < dims; i++ {
		if b.min[i] < r.min[i] || b.max[i] > r.max[i] {
			return false
		}
	}
	return true
}

func (r *box) largestAxis() (axis int, size float64) {
	j, jsz := 0, 0.0
	for i := 0; i < dims; i++ {
		sz := r.max[i] - r.min[i]
		if i == 0 || sz > jsz {
			j, jsz = i, sz
		}
	}
	return j, jsz
}

func (r *box) splitLargestAxisEdgeSnap(right *box) {
	axis, _ := r.largestAxis()
	left := r
	leftNode := left.data.(*node)
	rightNode := new(node)
	right.data = rightNode

	var equals []box
	for i := 0; i < leftNode.count; i++ {
		minDist := leftNode.boxes[i].min[axis] - left.min[axis]
		maxDist := left.max[axis] - leftNode.boxes[i].max[axis]
		if minDist < maxDist {
			// stay left
		} else {
			if minDist > maxDist {
				// move to right
				rightNode.boxes[rightNode.count] = leftNode.boxes[i]
				rightNode.count++
			} else {
				// move to equals, at the end of the left array
				equals = append(equals, leftNode.boxes[i])
			}
			leftNode.boxes[i] = leftNode.boxes[leftNode.count-1]
			leftNode.boxes[leftNode.count-1].data = nil
			leftNode.count--
			i--
		}
	}
	for _, b := range equals {
		if leftNode.count < rightNode.count {
			leftNode.boxes[leftNode.count] = b
			leftNode.count++
		} else {
			rightNode.boxes[rightNode.count] = b
			rightNode.count++
		}
	}
	left.recalc()
	right.recalc()
}

func (r *box) insert(item *box, height int) (grown bool) {
	n := r.data.(*node)
	if height == 0 {
		n.boxes[n.count] = *item
		n.count++
		grown = !r.contains(item)
		return grown
	}
	// choose subtree
	index := r.chooseLeastEnlargement(item)
	child := &n.boxes[index]
	grown = child.insert(item, height-1)
	if grown {
		child.expand(item)
		grown = !r.contains(item)
	}
	if child.data.(*node).count == maxEntries+1 {
		child.splitLargestAxisEdgeSnap(&n.boxes[n.count])
		n.count++
	}
	return grown
}

// fit an external item into a box type
func fit(min, max []float64, value interface{}, target *box) {
	if max == nil {
		max = min
	}
	if len(min) != len(max) {
		panic("min/max dimension mismatch")
	}
	if len(min) != dims {
		panic("invalid number of dimensions")
	}
	for i := 0; i < dims; i++ {
		target.min[i] = min[i]
		target.max[i] = max[i]
	}
	target.data = value
}

type overlapsResult int

const (
	not overlapsResult = iota
	intersects
	contains
)

// overlaps detects if r insersects or contains b.
// return not, intersects, contains
func (r *box) overlaps(b *box) overlapsResult {
	for i := 0; i < dims; i++ {
		if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
			return not
		}
		if r.min[i] > b.min[i] || b.max[i] > r.max[i] {
			i++
			for ; i < dims; i++ {
				if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
					return not
				}
			}
			return intersects
		}
	}
	return contains
}

// contains return struct when b is fully contained inside of n
func (r *box) intersects(b *box) bool {
	for i := 0; i < dims; i++ {
		if b.min[i] > r.max[i] || b.max[i] < r.min[i] {
			return false
		}
	}
	return true
}

func (r *box) search(
	target *box, height int,
	iter func(min, max []float64, value interface{}) bool,
) bool {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if target.intersects(&n.boxes[i]) {
				if !iter(n.boxes[i].min[:], n.boxes[i].max[:],
					n.boxes[i].data) {
					return false
				}
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			switch target.overlaps(&n.boxes[i]) {
			case intersects:
				if !n.boxes[i].search(target, height-1, iter) {
					return false
				}
			case contains:
				if !n.boxes[i].scan(target, height-1, iter) {
					return false
				}
			}
		}
	}
	return true
}

func (tr *BoxTree) search(
	target *box,
	iter func(min, max []float64, value interface{}) bool,
) {
	if tr.root.data == nil {
		return
	}
	res := target.overlaps(&tr.root)
	if res == intersects {
		tr.root.search(target, tr.height, iter)
	} else if res == contains {
		tr.root.scan(target, tr.height, iter)
	}
}

// Search ...
func (tr *BoxTree) Search(min, max []float64,
	iter func(min, max []float64, value interface{}) bool,
) {
	var target box
	fit(min, max, nil, &target)
	tr.search(&target, iter)
}

const (
	// Continue to first child box and/or next sibling.
	Continue = iota
	// Ignore child boxes but continue to next sibling.
	Ignore
	// Stop iterating
	Stop
)

// Traverse iterates through all items and container boxes in tree.
func (tr *BoxTree) Traverse(
	iter func(min, max []float64, height, level int, value interface{}) int,
) {
	if tr.root.data == nil {
		return
	}
	if iter(tr.root.min[:], tr.root.max[:], tr.height+1, 0, nil) == Continue {
		tr.root.traverse(tr.height, 1, iter)
	}
}

func (r *box) traverse(
	height, level int,
	iter func(min, max []float64, height, level int, value interface{}) int,
) int {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			action := iter(n.boxes[i].min[:], n.boxes[i].max[:], height, level,
				n.boxes[i].data)
			if action == Stop {
				return Stop
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			switch iter(n.boxes[i].min[:], n.boxes[i].max[:], height, level,
				n.boxes[i].data) {
			case Ignore:
			case Continue:
				if n.boxes[i].traverse(height-1, level+1, iter) == Stop {
					return Stop
				}
			case Stop:
				return Stop
			}
		}
	}
	return Continue
}

func (r *box) scan(
	target *box, height int,
	iter func(min, max []float64, value interface{}) bool,
) bool {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if !iter(n.boxes[i].min[:], n.boxes[i].max[:], n.boxes[i].data) {
				return false
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if !n.boxes[i].scan(target, height-1, iter) {
				return false
			}
		}
	}
	return true
}

// Scan iterates through all items in tree.
func (tr *BoxTree) Scan(iter func(min, max []float64, value interface{}) bool) {
	if tr.root.data == nil {
		return
	}
	tr.root.scan(nil, tr.height, iter)
}

// Delete ...
func (tr *BoxTree) Delete(min, max []float64, value interface{}) {
	var item box
	fit(min, max, value, &item)
	if tr.root.data == nil || !tr.root.contains(&item) {
		return
	}
	var removed, recalced bool
	removed, recalced, tr.reinsert =
		tr.root.delete(&item, tr.height, tr.reinsert[:0])
	if !removed {
		return
	}
	tr.count -= len(tr.reinsert) + 1
	if tr.count == 0 {
		tr.root = box{}
		recalced = false
	} else {
		for tr.height > 0 && tr.root.data.(*node).count == 1 {
			tr.root = tr.root.data.(*node).boxes[0]
			tr.height--
			tr.root.recalc()
		}
	}
	if recalced {
		tr.root.recalc()
	}
	for i := range tr.reinsert {
		tr.insert(&tr.reinsert[i])
		tr.reinsert[i].data = nil
	}
}

func (r *box) delete(item *box, height int, reinsert []box) (
	removed, recalced bool, reinsertOut []box,
) {
	n := r.data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if n.boxes[i].data == item.data {
				// found the target item to delete
				recalced = r.onEdge(&n.boxes[i])
				n.boxes[i] = n.boxes[n.count-1]
				n.boxes[n.count-1].data = nil
				n.count--
				if recalced {
					r.recalc()
				}
				return true, recalced, reinsert
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if !n.boxes[i].contains(item) {
				continue
			}
			removed, recalced, reinsert =
				n.boxes[i].delete(item, height-1, reinsert)
			if !removed {
				continue
			}
			if n.boxes[i].data.(*node).count < minEntries {
				// underflow
				if !recalced {
					recalced = r.onEdge(&n.boxes[i])
				}
				reinsert = n.boxes[i].flatten(reinsert, height-1)
				n.boxes[i] = n.boxes[n.count-1]
				n.boxes[n.count-1].data = nil
				n.count--
			}
			if recalced {
				r.recalc()
			}
			return removed, recalced, reinsert
		}
	}
	return false, false, reinsert
}

// flatten flattens all leaf boxes into a single list
func (r *box) flatten(all []box, height int) []box {
	n := r.data.(*node)
	if height == 0 {
		all = append(all, n.boxes[:n.count]...)
	} else {
		for i := 0; i < n.count; i++ {
			all = n.boxes[i].flatten(all, height-1)
		}
	}
	return all
}

// onedge returns true when b is on the edge of r
func (r *box) onEdge(b *box) bool {
	for i := 0; i < dims; i++ {
		if r.min[i] == b.min[i] || r.max[i] == b.max[i] {
			return true
		}
	}
	return false
}

// Count ...
func (tr *BoxTree) Count() int {
	return tr.count
}

func (r *box) totalOverlapArea(height int) float64 {
	var area float64
	n := r.data.(*node)
	for i := 0; i < n.count; i++ {
		for j := i + 1; j < n.count; j++ {
			area += n.boxes[i].overlapArea(&n.boxes[j])
		}

	}
	if height > 0 {
		for i := 0; i < n.count; i++ {
			area += n.boxes[i].totalOverlapArea(height - 1)
		}
	}
	return area
}

// TotalOverlapArea ...
func (tr *BoxTree) TotalOverlapArea() float64 {
	if tr.root.data == nil {
		return 0
	}
	return tr.root.totalOverlapArea(tr.height)
}

type qnode struct {
	dist float64
	box  box
}

type queue struct {
	nodes []qnode
	len   int
	size  int
}

func (q *queue) push(dist float64, box box) {
	if q.nodes == nil {
		q.nodes = make([]qnode, 2)
	} else {
		q.nodes = append(q.nodes, qnode{})
	}
	i := q.len + 1
	j := i / 2
	for i > 1 && q.nodes[j].dist > dist {
		q.nodes[i] = q.nodes[j]
		i = j
		j = j / 2
	}
	q.nodes[i].dist = dist
	q.nodes[i].box = box
	q.len++
}

func (q *queue) peek() qnode {
	if q.len == 0 {
		return qnode{}
	}
	return q.nodes[1]
}

func (q *queue) pop() qnode {
	if q.len == 0 {
		return qnode{}
	}
	n := q.nodes[1]
	q.nodes[1] = q.nodes[q.len]
	q.len--
	var j, k int
	i := 1
	for i != q.len+1 {
		k = q.len + 1
		j = 2 * i
		if j <= q.len && q.nodes[j].dist < q.nodes[k].dist {
			k = j
		}
		if j+1 <= q.len && q.nodes[j+1].dist < q.nodes[k].dist {
			k = j + 1
		}
		q.nodes[i] = q.nodes[k]
		i = k
	}
	return n
}

// Nearby returns items nearest to farthest.
// The dist param is the "box distance".
func (tr *BoxTree) Nearby(min, max []float64,
	iter func(min, max []float64, item interface{}) bool) {
	if tr.root.data == nil {
		return
	}
	var bbox box
	fit(min, max, nil, &bbox)
	box := tr.root
	var q queue
	for {
		n := box.data.(*node)
		for i := 0; i < n.count; i++ {
			dist := boxDist(&bbox, &n.boxes[i])
			q.push(dist, n.boxes[i])
		}
		for q.len > 0 {
			if _, ok := q.peek().box.data.(*node); ok {
				break
			}
			item := q.pop()
			if !iter(item.box.min[:], item.box.max[:], item.box.data) {
				return
			}
		}
		if q.len == 0 {
			break
		} else {
			box = q.pop().box
		}
	}
	return
}

func boxDist(a, b *box) float64 {
	var dist float64
	for i := 0; i < len(a.min); i++ {
		var min, max float64
		if a.min[i] > b.min[i] {
			min = a.min[i]
		} else {
			min = b.min[i]
		}
		if a.max[i] < b.max[i] {
			max = a.max[i]
		} else {
			max = b.max[i]
		}
		squared := min - max
		if squared > 0 {
			dist += squared * squared
		}
	}
	return dist
}

// Bounds returns the minimum bounding box
func (tr *BoxTree) Bounds() (min, max []float64) {
	if tr.root.data == nil {
		return
	}
	return tr.root.min[:], tr.root.max[:]
}
