package accel

import (
	"goraytracer/geometry"
	"goraytracer/material"
	"goraytracer/mesh"
	"goraytracer/ray"
	"goraytracer/vec3"
)

// TODO: rename this file to octtree.go
const MaxDepth = 5

type IntersectCandidate struct {
	geometry *geometry.Geometry
	material *material.Material
}

type OctTreeNode struct {
	aabb                geometry.AABB
	children            []OctTreeNode
	intersectCandidates []IntersectCandidate
}

type OctTree struct {
	rootNode OctTreeNode
}

func (node *OctTreeNode) search(ray *ray.Ray) []IntersectCandidate {
	if node.aabb.IntersectsRay(ray) == false {
		return nil
	}

	if node.children == nil && node.intersectCandidates != nil {
		return node.intersectCandidates
	}

	// recursively test all children
	for _, child := range node.children {
		childQuery := child.search(ray)
		if childQuery != nil {
			return childQuery
		}
	}

	return nil
}

func (tree *OctTree) Search(ray *ray.Ray) []IntersectCandidate {
	return tree.rootNode.search(ray)
}

func (tree *OctTree) Print() {
	// look at cockroach sql explain formatting for inspiration here
}

func geometryInAABB(meshes []mesh.Mesh, aabb geometry.AABB) []IntersectCandidate {
	// iterate over all meshes, and determine if mesh.geometry falls within aabb
	// I should also return the reference to the geometry's mesh
	// so I can store the reference to the geometry's mesh in the octree leaf node.

	candidates := make([]IntersectCandidate, 0)

	for _, mesh := range meshes {
		// a mesh could have more than one piece of geometry associated
		// a sphere mesh only has 1 sphere geometry
		// a polygonal mesh could have multiple triangles

		intersected := mesh.Geometry.AABBIntersections(aabb)

		for _, geometry := range intersected {
			candidates = append(candidates, IntersectCandidate{
				geometry: geometry,
				material: &mesh.Material,
			})
		}

	}

	return candidates
}

// split an AABB into 8 sub AABBs
func SplitAABB(parent geometry.AABB) []geometry.AABB {

	half := vec3.MultiplyScalar(vec3.Sub(parent.Max, parent.Min), 0.5)

	child0 := geometry.AABB{
		Min: parent.Min,
		Max: vec3.Sub(parent.Max, half)}

	child1 := geometry.AABB{
		Min: vec3.Add(parent.Min, vec3.Vec3{Z: half.Z}),
		Max: vec3.Sub(parent.Max, vec3.Vec3{X: half.X, Y: half.Y})}

	child2 := geometry.AABB{
		Min: vec3.Add(parent.Min, vec3.Vec3{X: half.X, Z: half.Z}),
		Max: vec3.Sub(parent.Max, vec3.Vec3{Y: half.Y})}

	child3 := geometry.AABB{
		Min: vec3.Add(parent.Min, vec3.Vec3{X: half.X}),
		Max: vec3.Sub(parent.Max, vec3.Vec3{Y: half.Y, Z: half.Z})}

	child4 := geometry.AABB{
		Min: vec3.Add(parent.Min, vec3.Vec3{Y: half.Y}),
		Max: vec3.Sub(parent.Max, vec3.Vec3{X: half.X, Z: half.Z})}

	child5 := geometry.AABB{
		Min: vec3.Add(parent.Min, vec3.Vec3{Y: half.Y, Z: half.Z}),
		Max: vec3.Sub(parent.Max, vec3.Vec3{X: half.X})}

	child6 := geometry.AABB{
		Min: vec3.Add(parent.Min, half),
		Max: parent.Max}

	child7 := geometry.AABB{
		Min: vec3.Add(parent.Min, vec3.Vec3{X: half.X, Y: half.Y}),
		Max: vec3.Sub(parent.Max, vec3.Vec3{Z: half.Z})}

	children := make([]geometry.AABB, 8)
	children[0] = child0
	children[1] = child1
	children[2] = child2
	children[3] = child3
	children[4] = child4
	children[5] = child5
	children[6] = child6
	children[7] = child7
	return children
}

func buildOctTreeNode(meshes []mesh.Mesh, aabb geometry.AABB, depth int) OctTreeNode {
	node := OctTreeNode{
		aabb:                aabb,
		children:            nil,
		intersectCandidates: nil,
	}

	if geo := geometryInAABB(meshes, aabb); len(geo) > 0 {
		if depth < MaxDepth {
			children := make([]OctTreeNode, 8)
			childAABBs := SplitAABB(aabb)
			for i, childAABB := range childAABBs {
				children[i] = buildOctTreeNode(meshes, childAABB, depth+1)
			}
		} else {
			node.intersectCandidates = geo
		}
	}

	return node
}

func BuildOctTree(meshes []mesh.Mesh) OctTree {
	// assume world extends from {-1,-1,-1} to {1,1,1}
	return OctTree{buildOctTreeNode(meshes, geometry.AABB{
		Min: vec3.Vec3{X: -1, Y: -1, Z: -1},
		Max: vec3.Vec3{X: 1, Y: 1, Z: 1},
	}, 0)}
}
