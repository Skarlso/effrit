package pkg

// Analyse will calculate stability and abstraction metrics
// for each package discovered by Scan.
func Analyse(packages *Packages) *Packages {
	for k, v := range packages.packageMap {
		v.Stability = calculateStability(v)
		packages.packageMap[k] = v
	}
	return packages
}

func calculateStability(p Package) float64 {
	return p.ImportCount / (p.ImportCount + p.DependedOnByCount)
}

// TODO: use go AST to count structs and interfaces ratio
func calculateAbstractness(p Package) float64 {
	return 0.0
}

func calculateDistanceFromMainSequence(p Package) float64 {
	return 0.0
}
