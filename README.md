![Ifrit](./img/ifrit.png)

# Effrit

[![Build Status](https://travis-ci.org/Skarlso/effrit.svg?branch=master)](https://travis-ci.org/Skarlso/effrit)

Go Efferent and Afferent package metric calculator.

Metrics calculated currently by this package:

- [x] Instability metric
- [x] Abstractness metric
- [x] Distance from main sequence metric

All metrics are now supported and calculated accordingly.

# Description of these metrics

https://en.wikipedia.org/wiki/Software_package_metrics

Please see Robert Cecil Martin's Clean Architecture book on details describing those metrics.

In terms of what this tool is doing, please refer to this post: [Efferent and Afferent Metrics in Go](https://skarlso.github.io/2019/04/21/efferent-and-afferent-metrics-in-go/).

# Usage on Effrit

Here is an example of running this tool on this very project:

![screenshot](./img/effrit_package.png)

# Package Data

Effrit now dumps data as JSON file into the project root directly. Until I finish the CGUI for effrit, this data can be processed by any other tool. Example using Effrit project:

```json
{
   "packages":[
      {
         "Name":"effrit",
         "FullName":"github.com/Skarlso/effrit",
         "Imports":[
            "github.com/Skarlso/effrit/cmd"
         ],
         "ImportCount":1,
         "DependedOnByCount":0,
         "DependedOnByNames":null,
         "Stability":1,
         "Abstractness":0,
         "DistanceFromMedian":0,
         "Dir":"/Users/hannibal/goprojects/effrit",
         "GoFiles":[
            "main.go"
         ]
      },
      {
         "Name":"cmd",
         "FullName":"github.com/Skarlso/effrit/cmd",
         "Imports":[
            "github.com/Skarlso/effrit/pkg"
         ],
         "ImportCount":1,
         "DependedOnByCount":1,
         "DependedOnByNames":[
            "github.com/Skarlso/effrit"
         ],
         "Stability":0.5,
         "Abstractness":0.5,
         "DistanceFromMedian":0,
         "Dir":"/Users/hannibal/goprojects/effrit/cmd",
         "GoFiles":[
            "root.go",
            "scan.go"
         ]
      },
      {
         "Name":"pkg",
         "FullName":"github.com/Skarlso/effrit/pkg",
         "Imports":[

         ],
         "ImportCount":0,
         "DependedOnByCount":1,
         "DependedOnByNames":[
            "github.com/Skarlso/effrit/cmd"
         ],
         "Stability":0,
         "Abstractness":0.3,
         "DistanceFromMedian":0.7,
         "Dir":"/Users/hannibal/goprojects/effrit/pkg",
         "GoFiles":[
            "packages.go",
            "scan.go"
         ]
      }
   ]
}
```

# Contributions

Are always welcomed!
