# Collections reconciler (ko-rek)

[![codecov](https://codecov.io/gh/rkosegi/go-korek/graph/badge.svg?token=BG1D2QKXRE)](https://codecov.io/gh/rkosegi/go-korek)
[![Go Report Card](https://goreportcard.com/badge/github.com/rkosegi/go-korek)](https://goreportcard.com/report/github.com/rkosegi/go-korek)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=rkosegi_go-korek&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=rkosegi_go-korek)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=rkosegi_go-korek&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=rkosegi_go-korek)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=rkosegi_go-korek&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=rkosegi_go-korek)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=rkosegi_go-korek&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=rkosegi_go-korek)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=rkosegi_go-korek&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=rkosegi_go-korek)
[![Go Reference](https://pkg.go.dev/badge/github.com/rkosegi/go-korek.svg)](https://pkg.go.dev/github.com/rkosegi/go-korek)
[![Apache 2.0 License](https://badgen.net/static/license/Apache2.0/blue)](https://github.com/rkosegi/go-korek/blob/main/LICENSE)
[![CodeQL Status](https://github.com/rkosegi/go-korek/actions/workflows/codeql.yaml/badge.svg)](https://github.com/rkosegi/go-korek/security/code-scanning)
[![CI Status](https://github.com/rkosegi/go-korek/actions/workflows/ci.yaml/badge.svg)](https://github.com/rkosegi/go-korek/actions/workflows/ci.yaml)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/rkosegi/go-korek/badge)](https://scorecard.dev/viewer/?uri=github.com/rkosegi/go-korek)


This is a library to simplify reconciliation of collections, in Go.
Common use-case is to reconcile system state from 2 different sources, such as external API and internal data structures.
This library can compute difference between 2 collections (slices or maps) and gives you back information about what was changed,
what's new and what was removed.

## Example

You can check unit test that tests similar thing [here](slice_test.go).

```go
type emp struct {
    Name       string
    Department string
    Salary     int
}

sr := ForSlice[emp]().
    WithEqualityFunc(DefaultEqualityFunc[emp]()).
    WithIdentityFunc(func(left, right emp) bool {
        return left.Name == right.Name
    })

left := []emp{
    {Name: "Alice", Department: "HR", Salary: 10},
    {Name: "Bob", Department: "IT", Salary: 10},
    {Name: "Charlie", Department: "Toilets", Salary: 99},
}
right := []emp{
    {Name: "Alice", Department: "HR", Salary: 10},
    {Name: "Bob", Department: "IT", Salary: 20},
    {Name: "Cyril", Department: "Sales", Salary: 10},
    {Name: "Dave", Department: "Management", Salary: 1},
}
same, changed, onlyLeft, onlyRight := sr.Diff(left, right)

fmt.Printf("Same: %v\n", same)
fmt.Printf("Changed: %v\n", changed)
fmt.Printf("Only in left: %v\n", onlyLeft)
fmt.Printf("Only in right: %v\n", onlyRight)
```

Now you can act on changed items according to your business logic, e.g. create items in internal data structure based on what is missing or the opposite,
remove items that are not present.
