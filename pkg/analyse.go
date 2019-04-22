package pkg

// Analyse will calculate stability and abstraction metrics
// for each package discovered by Scan.
func Analyse(packages map[string]Package) map[string]Package {
	for k, v := range packages {
		v.Stability = calculateStability(v)
		packages[k] = v
	}
	return packages
}

func calculateStability(p Package) float64 {
	return p.ImportCount / (p.ImportCount + p.DependedOnByCount)
}
