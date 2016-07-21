// generated; DO NOT EDIT!

/*

TITLE

	R-TREES: A DYNAMIC INDEX STRUCTURE FOR SPATIAL SEARCHING

DESCRIPTION

	A Go version of the RTree algorithm.

AUTHORS

	* 1983 Original algorithm and test code by Antonin Guttman and Michael Stonebraker, UC Berkely
	* 1994 ANCI C ported from original test code by Melinda Green - melinda@superliminal.com
	* 1995 Sphere volume fix for degeneracy problem submitted by Paul Brook
	* 2004 Templated C++ port by Greg Douglas
	* 2016 Go port by Josh Baker

LICENSE:

	Entirely free for all uses. Enjoy!

*/

// Implementation of RTree, a multidimensional bounding rectangle tree.
package rtree

import "math"

func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func ASSERT(condition bool) {
	if _DEBUG && !condition {
		panic("assertion failed")
	}
}

const (
	_DEBUG               = false
	NUMDIMS              = 18
	MAXNODES             = 8
	MINNODES             = MAXNODES / 2
	USE_SPHERICAL_VOLUME = true // Better split classification, may be slower on some systems
)

type ResultCallback func(dataID interface{}) bool

var unitSphereVolume = []float64{
	0.000000, 2.000000, 3.141593, // Dimension  0,1,2
	4.188790, 4.934802, 5.263789, // Dimension  3,4,5
	5.167713, 4.724766, 4.058712, // Dimension  6,7,8
	3.298509, 2.550164, 1.884104, // Dimension  9,10,11
	1.335263, 0.910629, 0.599265, // Dimension  12,13,14
	0.381443, 0.235331, 0.140981, // Dimension  15,16,17
	0.082146, 0.046622, 0.025807, // Dimension  18,19,20
}[NUMDIMS]

type RTree struct {
	root *Node ///< Root of tree
}

/// Minimal bounding rectangle (n-dimensional)
type Rect struct {
	min [NUMDIMS]float64 ///< Min dimensions of bounding box
	max [NUMDIMS]float64 ///< Max dimensions of bounding box
}

/// May be data or may be another subtree
/// The parents level determines this.
/// If the parents level is 0, then this is data
type Branch struct {
	rect  Rect        ///< Bounds
	child *Node       ///< Child node
	data  interface{} ///< Data Id or Ptr
}

/// Node for each branch level
type Node struct {
	count  int              ///< Count
	level  int              ///< Leaf is zero, others positive
	branch [MAXNODES]Branch ///< Branch
}

func (node *Node) IsInternalNode() bool {
	return (node.level > 0) // Not a leaf, but a internal node
}
func (node *Node) IsLeaf() bool {
	return (node.level == 0) // A leaf, contains data
}

/// A link list of nodes for reinsertion after a delete operation
type ListNode struct {
	next *ListNode ///< Next in list
	node *Node     ///< Node
}

const NOT_TAKEN = -1 // indicates that position

/// Variables for finding a split partition
type PartitionVars struct {
	partition [MAXNODES + 1]int
	total     int
	minFill   int
	count     [2]int
	cover     [2]Rect
	area      [2]float64

	branchBuf      [MAXNODES + 1]Branch
	branchCount    int
	coverSplit     Rect
	coverSplitArea float64
}

func NewRTree() *RTree {
	ASSERT(MAXNODES > MINNODES)
	ASSERT(MINNODES > 0)

	// We only support machine word size simple data type eg. integer index or object pointer.
	// Since we are storing as union with non data branch
	//ASSERT(sizeof(DATATYPE) == sizeof(void*) || sizeof(DATATYPE) == sizeof(int));

	return &RTree{
		root: &Node{},
	}
}

/// Insert entry
/// \param a_min Min of bounding rect
/// \param a_max Max of bounding rect
/// \param a_dataId Positive Id of data.  Maybe zero, but negative numbers not allowed.
func (tr *RTree) Insert(min, max [NUMDIMS]float64, dataId interface{}) {
	if _DEBUG {
		for index := 0; index < NUMDIMS; index++ {
			ASSERT(min[index] <= max[index])
		}
	} //_DEBUG
	var branch Branch
	branch.data = dataId
	for axis := 0; axis < NUMDIMS; axis++ {
		branch.rect.min[axis] = min[axis]
		branch.rect.max[axis] = max[axis]
	}
	InsertRect(&branch, &tr.root, 0)
}

/// Remove entry
/// \param a_min Min of bounding rect
/// \param a_max Max of bounding rect
/// \param a_dataId Positive Id of data.  Maybe zero, but negative numbers not allowed.
func (tr *RTree) Remove(min, max [NUMDIMS]float64, dataId interface{}) {
	if _DEBUG {
		for index := 0; index < NUMDIMS; index++ {
			ASSERT(min[index] <= max[index])
		}
	} //_DEBUG
	var rect Rect
	for axis := 0; axis < NUMDIMS; axis++ {
		rect.min[axis] = min[axis]
		rect.max[axis] = max[axis]
	}
	RemoveRect(&rect, dataId, &tr.root)
}

/// Find all within search rectangle
/// \param a_min Min of search bounding rect
/// \param a_max Max of search bounding rect
/// \param a_searchResult Search result array.  Caller should set grow size. Function will reset, not append to array.
/// \param a_resultCallback Callback function to return result.  Callback should return 'true' to continue searching
/// \param a_context User context to pass as parameter to a_resultCallback
/// \return Returns the number of entries found
func (tr *RTree) Search(min, max [NUMDIMS]float64, resultCallback ResultCallback) int {
	if _DEBUG {
		for index := 0; index < NUMDIMS; index++ {
			ASSERT(min[index] <= max[index])
		}
	} //_DEBUG
	var rect Rect
	for axis := 0; axis < NUMDIMS; axis++ {
		rect.min[axis] = min[axis]
		rect.max[axis] = max[axis]
	}
	foundCount, _ := Search(tr.root, rect, 0, resultCallback)
	return foundCount
}

/// Count the data elements in this container.  This is slow as no internal counter is maintained.
func (tr *RTree) Count() int {
	var count int
	CountRec(tr.root, &count)
	return count
}

func CountRec(node *Node, count *int) {
	if node.IsInternalNode() { // not a leaf node
		for index := 0; index < node.count; index++ {
			CountRec(node.branch[index].child, count)
		}
	} else { // A leaf node
		*count += node.count
	}
}

/// Remove all entries from tree
func (tr *RTree) RemoveAll() {
	// Delete all existing nodes
	tr.root = &Node{}
}

func InitNode(node *Node) {
	node.count = 0
	node.level = -1
}

func InitRect(rect *Rect) {
	for index := 0; index < NUMDIMS; index++ {
		rect.min[index] = 0
		rect.max[index] = 0
	}
}

// Inserts a new data rectangle into the index structure.
// Recursively descends tree, propagates splits back up.
// Returns 0 if node was not split.  Old node updated.
// If node was split, returns 1 and sets the pointer pointed to by
// new_node to point to the new node.  Old node updated to become one of two.
// The level argument specifies the number of steps up from the leaf
// level to insert; e.g. a data rectangle goes in at level = 0.
func InsertRectRec(branch *Branch, node *Node, newNode **Node, level int) bool {
	ASSERT(node != nil && newNode != nil)
	ASSERT(level >= 0 && level <= node.level)

	// recurse until we reach the correct level for the new record. data records
	// will always be called with a_level == 0 (leaf)
	if node.level > level {
		// Still above level for insertion, go down tree recursively
		var otherNode *Node
		//var newBranch Branch

		// find the optimal branch for this record
		index := PickBranch(&branch.rect, node)

		// recursively insert this record into the picked branch
		childWasSplit := InsertRectRec(branch, node.branch[index].child, &otherNode, level)

		if !childWasSplit {
			// Child was not split. Merge the bounding box of the new record with the
			// existing bounding box
			node.branch[index].rect = CombineRect(&branch.rect, &(node.branch[index].rect))
			return false
		} else {
			// Child was split. The old branches are now re-partitioned to two nodes
			// so we have to re-calculate the bounding boxes of each node
			node.branch[index].rect = NodeCover(node.branch[index].child)
			var newBranch Branch
			newBranch.child = otherNode
			newBranch.rect = NodeCover(otherNode)

			// The old node is already a child of a_node. Now add the newly-created
			// node to a_node as well. a_node might be split because of that.
			return AddBranch(&newBranch, node, newNode)
		}
	} else if node.level == level {
		// We have reached level for insertion. Add rect, split if necessary
		return AddBranch(branch, node, newNode)
	} else {
		// Should never occur
		ASSERT(false)
		return false
	}
}

// Insert a data rectangle into an index structure.
// InsertRect provides for splitting the root;
// returns 1 if root was split, 0 if it was not.
// The level argument specifies the number of steps up from the leaf
// level to insert; e.g. a data rectangle goes in at level = 0.
// InsertRect2 does the recursion.
//
func InsertRect(branch *Branch, root **Node, level int) bool {
	ASSERT(root != nil)
	ASSERT(level >= 0 && level <= (*root).level)
	if _DEBUG {
		for index := 0; index < NUMDIMS; index++ {
			ASSERT(branch.rect.min[index] <= branch.rect.max[index])
		}
	} //_DEBUG

	var newNode *Node

	if InsertRectRec(branch, *root, &newNode, level) { // Root split

		// Grow tree taller and new root
		newRoot := &Node{}
		newRoot.level = (*root).level + 1

		var newBranch Branch

		// add old root node as a child of the new root
		newBranch.rect = NodeCover(*root)
		newBranch.child = *root
		AddBranch(&newBranch, newRoot, nil)

		// add the split node as a child of the new root
		newBranch.rect = NodeCover(newNode)
		newBranch.child = newNode
		AddBranch(&newBranch, newRoot, nil)

		// set the new root as the root node
		*root = newRoot

		return true
	}

	return false
}

// Find the smallest rectangle that includes all rectangles in branches of a node.
func NodeCover(node *Node) Rect {
	ASSERT(node != nil)

	rect := node.branch[0].rect
	for index := 1; index < node.count; index++ {
		rect = CombineRect(&rect, &(node.branch[index].rect))
	}

	return rect
}

// Add a branch to a node.  Split the node if necessary.
// Returns 0 if node not split.  Old node updated.
// Returns 1 if node split, sets *new_node to address of new node.
// Old node updated, becomes one of two.
func AddBranch(branch *Branch, node *Node, newNode **Node) bool {
	ASSERT(branch != nil)
	ASSERT(node != nil)

	if node.count < MAXNODES { // Split won't be necessary

		node.branch[node.count] = *branch
		node.count++

		return false
	} else {
		ASSERT(newNode != nil)

		SplitNode(node, branch, newNode)
		return true
	}
}

// Disconnect a dependent node.
// Caller must return (or stop using iteration index) after this as count has changed
func DisconnectBranch(node *Node, index int) {
	ASSERT(node != nil && (index >= 0) && (index < MAXNODES))
	ASSERT(node.count > 0)

	// Remove element by swapping with the last element to prevent gaps in array
	node.branch[index] = node.branch[node.count-1]
	node.branch[node.count-1].data = nil
	node.branch[node.count-1].child = nil
	node.count--
}

// Pick a branch.  Pick the one that will need the smallest increase
// in area to accomodate the new rectangle.  This will result in the
// least total area for the covering rectangles in the current node.
// In case of a tie, pick the one which was smaller before, to get
// the best resolution when searching.
func PickBranch(rect *Rect, node *Node) int {
	ASSERT(rect != nil && node != nil)

	var firstTime bool = true
	var increase float64
	var bestIncr float64 = -1
	var area float64
	var bestArea float64
	var best int
	var tempRect Rect

	for index := 0; index < node.count; index++ {
		curRect := &node.branch[index].rect
		area = CalcRectVolume(curRect)
		tempRect = CombineRect(rect, curRect)
		increase = CalcRectVolume(&tempRect) - area
		if (increase < bestIncr) || firstTime {
			best = index
			bestArea = area
			bestIncr = increase
			firstTime = false
		} else if (increase == bestIncr) && (area < bestArea) {
			best = index
			bestArea = area
			bestIncr = increase
		}
	}
	return best
}

// Combine two rectangles into larger one containing both
func CombineRect(rectA, rectB *Rect) Rect {
	ASSERT(rectA != nil && rectB != nil)

	var newRect Rect

	for index := 0; index < NUMDIMS; index++ {
		newRect.min[index] = Min(rectA.min[index], rectB.min[index])
		newRect.max[index] = Max(rectA.max[index], rectB.max[index])
	}

	return newRect
}

// Split a node.
// Divides the nodes branches and the extra one between two nodes.
// Old node is one of the new ones, and one really new one is created.
// Tries more than one method for choosing a partition, uses best result.
func SplitNode(node *Node, branch *Branch, newNode **Node) {
	ASSERT(node != nil)
	ASSERT(branch != nil)

	// Could just use local here, but member or external is faster since it is reused
	var localVars PartitionVars
	parVars := &localVars

	// Load all the branches into a buffer, initialize old node
	GetBranches(node, branch, parVars)

	// Find partition
	ChoosePartition(parVars, MINNODES)

	// Create a new node to hold (about) half of the branches
	*newNode = &Node{}
	(*newNode).level = node.level

	// Put branches from buffer into 2 nodes according to the chosen partition
	node.count = 0
	LoadNodes(node, *newNode, parVars)

	ASSERT((node.count + (*newNode).count) == parVars.total)
}

// Calculate the n-dimensional volume of a rectangle
func RectVolume(rect *Rect) float64 {
	ASSERT(rect != nil)

	var volume float64 = 1

	for index := 0; index < NUMDIMS; index++ {
		volume *= rect.max[index] - rect.min[index]
	}

	ASSERT(volume >= 0)

	return volume
}

// The exact volume of the bounding sphere for the given Rect
func RectSphericalVolume(rect *Rect) float64 {
	ASSERT(rect != nil)

	var sumOfSquares float64 = 0
	var radius float64

	for index := 0; index < NUMDIMS; index++ {
		halfExtent := (rect.max[index] - rect.min[index]) * 0.5
		sumOfSquares += halfExtent * halfExtent
	}

	radius = math.Sqrt(sumOfSquares)

	// Pow maybe slow, so test for common dims just use x*x, x*x*x.
	if NUMDIMS == 5 {
		return (radius * radius * radius * radius * radius * unitSphereVolume)
	} else if NUMDIMS == 4 {
		return (radius * radius * radius * radius * unitSphereVolume)
	} else if NUMDIMS == 3 {
		return (radius * radius * radius * unitSphereVolume)
	} else if NUMDIMS == 2 {
		return (radius * radius * unitSphereVolume)
	} else {
		return (math.Pow(radius, NUMDIMS) * unitSphereVolume)
	}
}

// Use one of the methods to calculate retangle volume
func CalcRectVolume(rect *Rect) float64 {
	if USE_SPHERICAL_VOLUME {
		return RectSphericalVolume(rect) // Slower but helps certain merge cases
	} else { // RTREE_USE_SPHERICAL_VOLUME
		return RectVolume(rect) // Faster but can cause poor merges
	} // RTREE_USE_SPHERICAL_VOLUME
}

// Load branch buffer with branches from full node plus the extra branch.
func GetBranches(node *Node, branch *Branch, parVars *PartitionVars) {
	ASSERT(node != nil)
	ASSERT(branch != nil)

	ASSERT(node.count == MAXNODES)

	// Load the branch buffer
	for index := 0; index < MAXNODES; index++ {
		parVars.branchBuf[index] = node.branch[index]
	}
	parVars.branchBuf[MAXNODES] = *branch
	parVars.branchCount = MAXNODES + 1

	// Calculate rect containing all in the set
	parVars.coverSplit = parVars.branchBuf[0].rect
	for index := 1; index < MAXNODES+1; index++ {
		parVars.coverSplit = CombineRect(&parVars.coverSplit, &parVars.branchBuf[index].rect)
	}
	parVars.coverSplitArea = CalcRectVolume(&parVars.coverSplit)
}

// Method #0 for choosing a partition:
// As the seeds for the two groups, pick the two rects that would waste the
// most area if covered by a single rectangle, i.e. evidently the worst pair
// to have in the same group.
// Of the remaining, one at a time is chosen to be put in one of the two groups.
// The one chosen is the one with the greatest difference in area expansion
// depending on which group - the rect most strongly attracted to one group
// and repelled from the other.
// If one group gets too full (more would force other group to violate min
// fill requirement) then other group gets the rest.
// These last are the ones that can go in either group most easily.
func ChoosePartition(parVars *PartitionVars, minFill int) {
	ASSERT(parVars != nil)

	var biggestDiff float64
	var group, chosen, betterGroup int

	InitParVars(parVars, parVars.branchCount, minFill)
	PickSeeds(parVars)

	for ((parVars.count[0] + parVars.count[1]) < parVars.total) &&
		(parVars.count[0] < (parVars.total - parVars.minFill)) &&
		(parVars.count[1] < (parVars.total - parVars.minFill)) {
		biggestDiff = -1
		for index := 0; index < parVars.total; index++ {
			if NOT_TAKEN == parVars.partition[index] {
				curRect := &parVars.branchBuf[index].rect
				rect0 := CombineRect(curRect, &parVars.cover[0])
				rect1 := CombineRect(curRect, &parVars.cover[1])
				growth0 := CalcRectVolume(&rect0) - parVars.area[0]
				growth1 := CalcRectVolume(&rect1) - parVars.area[1]
				diff := growth1 - growth0
				if diff >= 0 {
					group = 0
				} else {
					group = 1
					diff = -diff
				}

				if diff > biggestDiff {
					biggestDiff = diff
					chosen = index
					betterGroup = group
				} else if (diff == biggestDiff) && (parVars.count[group] < parVars.count[betterGroup]) {
					chosen = index
					betterGroup = group
				}
			}
		}
		Classify(chosen, betterGroup, parVars)
	}

	// If one group too full, put remaining rects in the other
	if (parVars.count[0] + parVars.count[1]) < parVars.total {
		if parVars.count[0] >= parVars.total-parVars.minFill {
			group = 1
		} else {
			group = 0
		}
		for index := 0; index < parVars.total; index++ {
			if NOT_TAKEN == parVars.partition[index] {
				Classify(index, group, parVars)
			}
		}
	}

	ASSERT((parVars.count[0] + parVars.count[1]) == parVars.total)
	ASSERT((parVars.count[0] >= parVars.minFill) &&
		(parVars.count[1] >= parVars.minFill))
}

// Copy branches from the buffer into two nodes according to the partition.
func LoadNodes(nodeA, nodeB *Node, parVars *PartitionVars) {
	ASSERT(nodeA != nil)
	ASSERT(nodeB != nil)
	ASSERT(parVars != nil)

	for index := 0; index < parVars.total; index++ {
		ASSERT(parVars.partition[index] == 0 || parVars.partition[index] == 1)

		targetNodeIndex := parVars.partition[index]
		targetNodes := []*Node{nodeA, nodeB}

		// It is assured that AddBranch here will not cause a node split.
		nodeWasSplit := AddBranch(&parVars.branchBuf[index], targetNodes[targetNodeIndex], nil)
		ASSERT(!nodeWasSplit)
	}
}

// Initialize a PartitionVars structure.
func InitParVars(parVars *PartitionVars, maxRects, minFill int) {
	ASSERT(parVars != nil)

	parVars.count[0] = 0
	parVars.count[1] = 0
	parVars.area[0] = 0
	parVars.area[1] = 0
	parVars.total = maxRects
	parVars.minFill = minFill
	for index := 0; index < maxRects; index++ {
		parVars.partition[index] = NOT_TAKEN
	}
}

func PickSeeds(parVars *PartitionVars) {
	var seed0, seed1 int
	var worst, waste float64
	var area [MAXNODES + 1]float64

	for index := 0; index < parVars.total; index++ {
		area[index] = CalcRectVolume(&parVars.branchBuf[index].rect)
	}

	worst = -parVars.coverSplitArea - 1
	for indexA := 0; indexA < parVars.total-1; indexA++ {
		for indexB := indexA + 1; indexB < parVars.total; indexB++ {
			oneRect := CombineRect(&parVars.branchBuf[indexA].rect, &parVars.branchBuf[indexB].rect)
			waste = CalcRectVolume(&oneRect) - area[indexA] - area[indexB]
			if waste > worst {
				worst = waste
				seed0 = indexA
				seed1 = indexB
			}
		}
	}

	Classify(seed0, 0, parVars)
	Classify(seed1, 1, parVars)
}

// Put a branch in one of the groups.
func Classify(index, group int, parVars *PartitionVars) {
	ASSERT(parVars != nil)
	ASSERT(NOT_TAKEN == parVars.partition[index])

	parVars.partition[index] = group

	// Calculate combined rect
	if parVars.count[group] == 0 {
		parVars.cover[group] = parVars.branchBuf[index].rect
	} else {
		parVars.cover[group] = CombineRect(&parVars.branchBuf[index].rect, &parVars.cover[group])
	}

	// Calculate volume of combined rect
	parVars.area[group] = CalcRectVolume(&parVars.cover[group])

	parVars.count[group]++
}

// Delete a data rectangle from an index structure.
// Pass in a pointer to a Rect, the tid of the record, ptr to ptr to root node.
// Returns 1 if record not found, 0 if success.
// RemoveRect provides for eliminating the root.
func RemoveRect(rect *Rect, id interface{}, root **Node) bool {
	ASSERT(rect != nil && root != nil)
	ASSERT(*root != nil)

	var reInsertList *ListNode

	if !RemoveRectRec(rect, id, *root, &reInsertList) {
		// Found and deleted a data item
		// Reinsert any branches from eliminated nodes
		for reInsertList != nil {
			tempNode := reInsertList.node

			for index := 0; index < tempNode.count; index++ {
				// TODO go over this code. should I use (tempNode->m_level - 1)?
				InsertRect(&tempNode.branch[index], root, tempNode.level)
			}
			reInsertList = reInsertList.next
		}

		// Check for redundant root (not leaf, 1 child) and eliminate TODO replace
		// if with while? In case there is a whole branch of redundant roots...
		if (*root).count == 1 && (*root).IsInternalNode() {
			tempNode := (*root).branch[0].child

			ASSERT(tempNode != nil)
			*root = tempNode
		}
		return false
	} else {
		return true
	}
}

// Delete a rectangle from non-root part of an index structure.
// Called by RemoveRect.  Descends tree recursively,
// merges branches on the way back up.
// Returns 1 if record not found, 0 if success.
func RemoveRectRec(rect *Rect, id interface{}, node *Node, listNode **ListNode) bool {
	ASSERT(rect != nil && node != nil && listNode != nil)
	ASSERT(node.level >= 0)

	if node.IsInternalNode() { // not a leaf node
		for index := 0; index < node.count; index++ {
			if Overlap(*rect, node.branch[index].rect) {
				if !RemoveRectRec(rect, id, node.branch[index].child, listNode) {
					if node.branch[index].child.count >= MINNODES {
						// child removed, just resize parent rect
						node.branch[index].rect = NodeCover(node.branch[index].child)
					} else {
						// child removed, not enough entries in node, eliminate node
						ReInsert(node.branch[index].child, listNode)
						DisconnectBranch(node, index) // Must return after this call as count has changed
					}
					return false
				}
			}
		}
		return true
	} else { // A leaf node
		for index := 0; index < node.count; index++ {
			if node.branch[index].data == id {
				DisconnectBranch(node, index) // Must return after this call as count has changed
				return false
			}
		}
		return true
	}
}

// Decide whether two rectangles overlap.
func Overlap(rectA, rectB Rect) bool {
	for index := 0; index < NUMDIMS; index++ {
		if rectA.min[index] > rectB.max[index] ||
			rectB.min[index] > rectA.max[index] {
			return false
		}
	}
	return true
}

// Add a node to the reinsertion list.  All its branches will later
// be reinserted into the index structure.
func ReInsert(node *Node, listNode **ListNode) {
	newListNode := &ListNode{}
	newListNode.node = node
	newListNode.next = *listNode
	*listNode = newListNode
}

// Search in an index tree or subtree for all data retangles that overlap the argument rectangle.
func Search(node *Node, rect Rect, foundCount int, resultCallback ResultCallback) (int, bool) {
	ASSERT(node != nil)
	ASSERT(node.level >= 0)

	if node.IsInternalNode() {
		// This is an internal node in the tree
		for index := 0; index < node.count; index++ {
			if Overlap(rect, node.branch[index].rect) {
				var ok bool
				foundCount, ok = Search(node.branch[index].child, rect, foundCount, resultCallback)
				if !ok {
					// The callback indicated to stop searching
					return foundCount, false
				}
			}
		}
	} else {
		// This is a leaf node
		for index := 0; index < node.count; index++ {
			if Overlap(rect, node.branch[index].rect) {
				id := node.branch[index].data
				foundCount++
				if !resultCallback(id) {
					return foundCount, false // Don't continue searching
				}

			}
		}
	}
	return foundCount, true // Continue searching
}
