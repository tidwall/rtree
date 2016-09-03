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

func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

const (
	numDims            = 11
	maxNodes           = 8
	minNodes           = maxNodes / 2
	useSphericalVolume = true // Better split classification, may be slower on some systems
)

var unitSphereVolume = []float64{
	0.000000, 2.000000, 3.141593, // Dimension  0,1,2
	4.188790, 4.934802, 5.263789, // Dimension  3,4,5
	5.167713, 4.724766, 4.058712, // Dimension  6,7,8
	3.298509, 2.550164, 1.884104, // Dimension  9,10,11
	1.335263, 0.910629, 0.599265, // Dimension  12,13,14
	0.381443, 0.235331, 0.140981, // Dimension  15,16,17
	0.082146, 0.046622, 0.025807, // Dimension  18,19,20
}[numDims]

type RTree struct {
	root *nodeT ///< Root of tree
}

/// Minimal bounding rectangle (n-dimensional)
type rectT struct {
	min [numDims]float64 ///< Min dimensions of bounding box
	max [numDims]float64 ///< Max dimensions of bounding box
}

/// May be data or may be another subtree
/// The parents level determines this.
/// If the parents level is 0, then this is data
type branchT struct {
	rect  rectT       ///< Bounds
	child *nodeT      ///< Child node
	data  interface{} ///< Data Id or Ptr
}

/// nodeT for each branch level
type nodeT struct {
	count  int               ///< Count
	level  int               ///< Leaf is zero, others positive
	branch [maxNodes]branchT ///< Branch
}

func (node *nodeT) isInternalNode() bool {
	return (node.level > 0) // Not a leaf, but a internal node
}
func (node *nodeT) isLeaf() bool {
	return (node.level == 0) // A leaf, contains data
}

/// A link list of nodes for reinsertion after a delete operation
type listNodeT struct {
	next *listNodeT ///< Next in list
	node *nodeT     ///< Node
}

const notTaken = -1 // indicates that position

/// Variables for finding a split partition
type partitionVarsT struct {
	partition [maxNodes + 1]int
	total     int
	minFill   int
	count     [2]int
	cover     [2]rectT
	area      [2]float64

	branchBuf      [maxNodes + 1]branchT
	branchCount    int
	coverSplit     rectT
	coverSplitArea float64
}

func New() *RTree {
	// We only support machine word size simple data type eg. integer index or object pointer.
	// Since we are storing as union with non data branch
	return &RTree{
		root: &nodeT{},
	}
}

/// Insert entry
/// \param a_min Min of bounding rect
/// \param a_max Max of bounding rect
/// \param a_dataId Positive Id of data.  Maybe zero, but negative numbers not allowed.
func (tr *RTree) Insert(min, max [numDims]float64, dataId interface{}) {
	var branch branchT
	branch.data = dataId
	for axis := 0; axis < numDims; axis++ {
		branch.rect.min[axis] = min[axis]
		branch.rect.max[axis] = max[axis]
	}
	insertRect(&branch, &tr.root, 0)
}

/// Remove entry
/// \param a_min Min of bounding rect
/// \param a_max Max of bounding rect
/// \param a_dataId Positive Id of data.  Maybe zero, but negative numbers not allowed.
func (tr *RTree) Remove(min, max [numDims]float64, dataId interface{}) {
	var rect rectT
	for axis := 0; axis < numDims; axis++ {
		rect.min[axis] = min[axis]
		rect.max[axis] = max[axis]
	}
	removeRect(&rect, dataId, &tr.root)
}

/// Find all within search rectangle
/// \param a_min Min of search bounding rect
/// \param a_max Max of search bounding rect
/// \param a_searchResult search result array.  Caller should set grow size. Function will reset, not append to array.
/// \param a_resultCallback Callback function to return result.  Callback should return 'true' to continue searching
/// \param a_context User context to pass as parameter to a_resultCallback
/// \return Returns the number of entries found
func (tr *RTree) Search(min, max [numDims]float64, resultCallback func(data interface{}) bool) int {
	var rect rectT
	for axis := 0; axis < numDims; axis++ {
		rect.min[axis] = min[axis]
		rect.max[axis] = max[axis]
	}
	foundCount, _ := search(tr.root, rect, 0, resultCallback)
	return foundCount
}

/// Count the data elements in this container.  This is slow as no internal counter is maintained.
func (tr *RTree) Count() int {
	var count int
	countRec(tr.root, &count)
	return count
}

/// Remove all entries from tree
func (tr *RTree) RemoveAll() {
	// Delete all existing nodes
	tr.root = &nodeT{}
}

func countRec(node *nodeT, count *int) {
	if node.isInternalNode() { // not a leaf node
		for index := 0; index < node.count; index++ {
			countRec(node.branch[index].child, count)
		}
	} else { // A leaf node
		*count += node.count
	}
}

// Inserts a new data rectangle into the index structure.
// Recursively descends tree, propagates splits back up.
// Returns 0 if node was not split.  Old node updated.
// If node was split, returns 1 and sets the pointer pointed to by
// new_node to point to the new node.  Old node updated to become one of two.
// The level argument specifies the number of steps up from the leaf
// level to insert; e.g. a data rectangle goes in at level = 0.
func insertRectRec(branch *branchT, node *nodeT, newNode **nodeT, level int) bool {
	// recurse until we reach the correct level for the new record. data records
	// will always be called with a_level == 0 (leaf)
	if node.level > level {
		// Still above level for insertion, go down tree recursively
		var otherNode *nodeT
		//var newBranch branchT

		// find the optimal branch for this record
		index := pickBranch(&branch.rect, node)

		// recursively insert this record into the picked branch
		childWasSplit := insertRectRec(branch, node.branch[index].child, &otherNode, level)

		if !childWasSplit {
			// Child was not split. Merge the bounding box of the new record with the
			// existing bounding box
			node.branch[index].rect = combineRect(&branch.rect, &(node.branch[index].rect))
			return false
		} else {
			// Child was split. The old branches are now re-partitioned to two nodes
			// so we have to re-calculate the bounding boxes of each node
			node.branch[index].rect = nodeCover(node.branch[index].child)
			var newBranch branchT
			newBranch.child = otherNode
			newBranch.rect = nodeCover(otherNode)

			// The old node is already a child of a_node. Now add the newly-created
			// node to a_node as well. a_node might be split because of that.
			return addBranch(&newBranch, node, newNode)
		}
	} else if node.level == level {
		// We have reached level for insertion. Add rect, split if necessary
		return addBranch(branch, node, newNode)
	} else {
		// Should never occur
		return false
	}
}

// Insert a data rectangle into an index structure.
// insertRect provides for splitting the root;
// returns 1 if root was split, 0 if it was not.
// The level argument specifies the number of steps up from the leaf
// level to insert; e.g. a data rectangle goes in at level = 0.
// InsertRect2 does the recursion.
//
func insertRect(branch *branchT, root **nodeT, level int) bool {
	var newNode *nodeT

	if insertRectRec(branch, *root, &newNode, level) { // Root split

		// Grow tree taller and new root
		newRoot := &nodeT{}
		newRoot.level = (*root).level + 1

		var newBranch branchT

		// add old root node as a child of the new root
		newBranch.rect = nodeCover(*root)
		newBranch.child = *root
		addBranch(&newBranch, newRoot, nil)

		// add the split node as a child of the new root
		newBranch.rect = nodeCover(newNode)
		newBranch.child = newNode
		addBranch(&newBranch, newRoot, nil)

		// set the new root as the root node
		*root = newRoot

		return true
	}
	return false
}

// Find the smallest rectangle that includes all rectangles in branches of a node.
func nodeCover(node *nodeT) rectT {
	rect := node.branch[0].rect
	for index := 1; index < node.count; index++ {
		rect = combineRect(&rect, &(node.branch[index].rect))
	}
	return rect
}

// Add a branch to a node.  Split the node if necessary.
// Returns 0 if node not split.  Old node updated.
// Returns 1 if node split, sets *new_node to address of new node.
// Old node updated, becomes one of two.
func addBranch(branch *branchT, node *nodeT, newNode **nodeT) bool {
	if node.count < maxNodes { // Split won't be necessary
		node.branch[node.count] = *branch
		node.count++
		return false
	} else {
		splitNode(node, branch, newNode)
		return true
	}
}

// Disconnect a dependent node.
// Caller must return (or stop using iteration index) after this as count has changed
func disconnectBranch(node *nodeT, index int) {
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
func pickBranch(rect *rectT, node *nodeT) int {
	var firstTime bool = true
	var increase float64
	var bestIncr float64 = -1
	var area float64
	var bestArea float64
	var best int
	var tempRect rectT

	for index := 0; index < node.count; index++ {
		curRect := &node.branch[index].rect
		area = calcRectVolume(curRect)
		tempRect = combineRect(rect, curRect)
		increase = calcRectVolume(&tempRect) - area
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
func combineRect(rectA, rectB *rectT) rectT {
	var newRect rectT

	for index := 0; index < numDims; index++ {
		newRect.min[index] = fmin(rectA.min[index], rectB.min[index])
		newRect.max[index] = fmax(rectA.max[index], rectB.max[index])
	}

	return newRect
}

// Split a node.
// Divides the nodes branches and the extra one between two nodes.
// Old node is one of the new ones, and one really new one is created.
// Tries more than one method for choosing a partition, uses best result.
func splitNode(node *nodeT, branch *branchT, newNode **nodeT) {
	// Could just use local here, but member or external is faster since it is reused
	var localVars partitionVarsT
	parVars := &localVars

	// Load all the branches into a buffer, initialize old node
	getBranches(node, branch, parVars)

	// Find partition
	choosePartition(parVars, minNodes)

	// Create a new node to hold (about) half of the branches
	*newNode = &nodeT{}
	(*newNode).level = node.level

	// Put branches from buffer into 2 nodes according to the chosen partition
	node.count = 0
	loadNodes(node, *newNode, parVars)
}

// Calculate the n-dimensional volume of a rectangle
func rectVolume(rect *rectT) float64 {
	var volume float64 = 1
	for index := 0; index < numDims; index++ {
		volume *= rect.max[index] - rect.min[index]
	}
	return volume
}

// The exact volume of the bounding sphere for the given rectT
func rectSphericalVolume(rect *rectT) float64 {
	var sumOfSquares float64 = 0
	var radius float64

	for index := 0; index < numDims; index++ {
		halfExtent := (rect.max[index] - rect.min[index]) * 0.5
		sumOfSquares += halfExtent * halfExtent
	}

	radius = math.Sqrt(sumOfSquares)

	// Pow maybe slow, so test for common dims just use x*x, x*x*x.
	if numDims == 5 {
		return (radius * radius * radius * radius * radius * unitSphereVolume)
	} else if numDims == 4 {
		return (radius * radius * radius * radius * unitSphereVolume)
	} else if numDims == 3 {
		return (radius * radius * radius * unitSphereVolume)
	} else if numDims == 2 {
		return (radius * radius * unitSphereVolume)
	} else {
		return (math.Pow(radius, numDims) * unitSphereVolume)
	}
}

// Use one of the methods to calculate retangle volume
func calcRectVolume(rect *rectT) float64 {
	if useSphericalVolume {
		return rectSphericalVolume(rect) // Slower but helps certain merge cases
	} else { // RTREE_USE_SPHERICAL_VOLUME
		return rectVolume(rect) // Faster but can cause poor merges
	} // RTREE_USE_SPHERICAL_VOLUME
}

// Load branch buffer with branches from full node plus the extra branch.
func getBranches(node *nodeT, branch *branchT, parVars *partitionVarsT) {
	// Load the branch buffer
	for index := 0; index < maxNodes; index++ {
		parVars.branchBuf[index] = node.branch[index]
	}
	parVars.branchBuf[maxNodes] = *branch
	parVars.branchCount = maxNodes + 1

	// Calculate rect containing all in the set
	parVars.coverSplit = parVars.branchBuf[0].rect
	for index := 1; index < maxNodes+1; index++ {
		parVars.coverSplit = combineRect(&parVars.coverSplit, &parVars.branchBuf[index].rect)
	}
	parVars.coverSplitArea = calcRectVolume(&parVars.coverSplit)
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
func choosePartition(parVars *partitionVarsT, minFill int) {
	var biggestDiff float64
	var group, chosen, betterGroup int

	initParVars(parVars, parVars.branchCount, minFill)
	pickSeeds(parVars)

	for ((parVars.count[0] + parVars.count[1]) < parVars.total) &&
		(parVars.count[0] < (parVars.total - parVars.minFill)) &&
		(parVars.count[1] < (parVars.total - parVars.minFill)) {
		biggestDiff = -1
		for index := 0; index < parVars.total; index++ {
			if notTaken == parVars.partition[index] {
				curRect := &parVars.branchBuf[index].rect
				rect0 := combineRect(curRect, &parVars.cover[0])
				rect1 := combineRect(curRect, &parVars.cover[1])
				growth0 := calcRectVolume(&rect0) - parVars.area[0]
				growth1 := calcRectVolume(&rect1) - parVars.area[1]
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
		classify(chosen, betterGroup, parVars)
	}

	// If one group too full, put remaining rects in the other
	if (parVars.count[0] + parVars.count[1]) < parVars.total {
		if parVars.count[0] >= parVars.total-parVars.minFill {
			group = 1
		} else {
			group = 0
		}
		for index := 0; index < parVars.total; index++ {
			if notTaken == parVars.partition[index] {
				classify(index, group, parVars)
			}
		}
	}
}

// Copy branches from the buffer into two nodes according to the partition.
func loadNodes(nodeA, nodeB *nodeT, parVars *partitionVarsT) {
	for index := 0; index < parVars.total; index++ {
		targetNodeIndex := parVars.partition[index]
		targetNodes := []*nodeT{nodeA, nodeB}

		// It is assured that addBranch here will not cause a node split.
		addBranch(&parVars.branchBuf[index], targetNodes[targetNodeIndex], nil)
	}
}

// Initialize a partitionVarsT structure.
func initParVars(parVars *partitionVarsT, maxRects, minFill int) {
	parVars.count[0] = 0
	parVars.count[1] = 0
	parVars.area[0] = 0
	parVars.area[1] = 0
	parVars.total = maxRects
	parVars.minFill = minFill
	for index := 0; index < maxRects; index++ {
		parVars.partition[index] = notTaken
	}
}

func pickSeeds(parVars *partitionVarsT) {
	var seed0, seed1 int
	var worst, waste float64
	var area [maxNodes + 1]float64

	for index := 0; index < parVars.total; index++ {
		area[index] = calcRectVolume(&parVars.branchBuf[index].rect)
	}

	worst = -parVars.coverSplitArea - 1
	for indexA := 0; indexA < parVars.total-1; indexA++ {
		for indexB := indexA + 1; indexB < parVars.total; indexB++ {
			oneRect := combineRect(&parVars.branchBuf[indexA].rect, &parVars.branchBuf[indexB].rect)
			waste = calcRectVolume(&oneRect) - area[indexA] - area[indexB]
			if waste > worst {
				worst = waste
				seed0 = indexA
				seed1 = indexB
			}
		}
	}

	classify(seed0, 0, parVars)
	classify(seed1, 1, parVars)
}

// Put a branch in one of the groups.
func classify(index, group int, parVars *partitionVarsT) {
	parVars.partition[index] = group

	// Calculate combined rect
	if parVars.count[group] == 0 {
		parVars.cover[group] = parVars.branchBuf[index].rect
	} else {
		parVars.cover[group] = combineRect(&parVars.branchBuf[index].rect, &parVars.cover[group])
	}

	// Calculate volume of combined rect
	parVars.area[group] = calcRectVolume(&parVars.cover[group])

	parVars.count[group]++
}

// Delete a data rectangle from an index structure.
// Pass in a pointer to a rectT, the tid of the record, ptr to ptr to root node.
// Returns 1 if record not found, 0 if success.
// removeRect provides for eliminating the root.
func removeRect(rect *rectT, id interface{}, root **nodeT) bool {
	var reInsertList *listNodeT

	if !removeRectRec(rect, id, *root, &reInsertList) {
		// Found and deleted a data item
		// Reinsert any branches from eliminated nodes
		for reInsertList != nil {
			tempNode := reInsertList.node

			for index := 0; index < tempNode.count; index++ {
				// TODO go over this code. should I use (tempNode->m_level - 1)?
				insertRect(&tempNode.branch[index], root, tempNode.level)
			}
			reInsertList = reInsertList.next
		}

		// Check for redundant root (not leaf, 1 child) and eliminate TODO replace
		// if with while? In case there is a whole branch of redundant roots...
		if (*root).count == 1 && (*root).isInternalNode() {
			tempNode := (*root).branch[0].child
			*root = tempNode
		}
		return false
	} else {
		return true
	}
}

// Delete a rectangle from non-root part of an index structure.
// Called by removeRect.  Descends tree recursively,
// merges branches on the way back up.
// Returns 1 if record not found, 0 if success.
func removeRectRec(rect *rectT, id interface{}, node *nodeT, listNode **listNodeT) bool {
	if node.isInternalNode() { // not a leaf node
		for index := 0; index < node.count; index++ {
			if overlap(*rect, node.branch[index].rect) {
				if !removeRectRec(rect, id, node.branch[index].child, listNode) {
					if node.branch[index].child.count >= minNodes {
						// child removed, just resize parent rect
						node.branch[index].rect = nodeCover(node.branch[index].child)
					} else {
						// child removed, not enough entries in node, eliminate node
						reInsert(node.branch[index].child, listNode)
						disconnectBranch(node, index) // Must return after this call as count has changed
					}
					return false
				}
			}
		}
		return true
	} else { // A leaf node
		for index := 0; index < node.count; index++ {
			if node.branch[index].data == id {
				disconnectBranch(node, index) // Must return after this call as count has changed
				return false
			}
		}
		return true
	}
}

// Decide whether two rectangles overlap.
func overlap(rectA, rectB rectT) bool {
	for index := 0; index < numDims; index++ {
		if rectA.min[index] > rectB.max[index] ||
			rectB.min[index] > rectA.max[index] {
			return false
		}
	}
	return true
}

// Add a node to the reinsertion list.  All its branches will later
// be reinserted into the index structure.
func reInsert(node *nodeT, listNode **listNodeT) {
	newListNode := &listNodeT{}
	newListNode.node = node
	newListNode.next = *listNode
	*listNode = newListNode
}

// search in an index tree or subtree for all data retangles that overlap the argument rectangle.
func search(node *nodeT, rect rectT, foundCount int, resultCallback func(data interface{}) bool) (int, bool) {
	if node.isInternalNode() {
		// This is an internal node in the tree
		for index := 0; index < node.count; index++ {
			if overlap(rect, node.branch[index].rect) {
				var ok bool
				foundCount, ok = search(node.branch[index].child, rect, foundCount, resultCallback)
				if !ok {
					// The callback indicated to stop searching
					return foundCount, false
				}
			}
		}
	} else {
		// This is a leaf node
		for index := 0; index < node.count; index++ {
			if overlap(rect, node.branch[index].rect) {
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
