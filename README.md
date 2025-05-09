semver is a [Semantic Versioning](http://semver.org/) library written in golang. It fully covers spec version `2.0.0`.

This fork of [blang/semver](https://github.com/blang/semver) has been updated to accept 4 digit versions. This goes against
semver by incorporating a revision number which is used by many enterprise organizations.

Versioning
----------
Old v1-v3 versions exist in the root of the repository for compatiblity reasons and will only receive bug fixes.

The current stable version is [*v4*](v4/) and is fully go-mod compatible.

Usage
-----
```bash
$ go get github.com/khulnasoft/semver/v4
# Or use fixed versions
$ go get github.com/khulnasoft/semver/v4@v4.0.0
```
Note: Always vendor your dependencies or fix on a specific version tag.

```go
import github.com/khulnasoft/semver/v4
v1, err := semver.Make("1.0.0.0-beta")
v2, err := semver.Make("2.0.0.0-beta")
v1.Compare(v2)
```

Also check the [GoDocs](http://godoc.org/github.com/blang/semver/v4).

Why should I use this lib?
-----

- Fully spec compatible
- No reflection
- ~No regex~
  - Regex was added to make parsing 3 and 4 digit versioning for readability of code
- Fully tested (Coverage >99%)
- Readable parsing/validation errors
- Fast (See [Benchmarks](#benchmarks))
- Only Stdlib
- Uses values instead of pointers
- Many features, see below


Features
-----

- Parsing and validation at all levels
- Comparator-like comparisons
- Compare Helper Methods
- InPlace manipulation
- Ranges `>=1.0.0 <2.0.0 || >=3.0.0 !3.0.1-beta.1`, `>=1.0.0.0 <2.0.0.0 || >=3.0.0 !3.0.1.0-beta.1`
- Wildcards `>=1.x`, `<=2.5.x`, `>=1.2.3.x`, `<=2.5.x`
- Sortable (implements sort.Interface)
- database/sql compatible (sql.Scanner/Valuer)
- encoding/json compatible (json.Marshaler/Unmarshaler)

Ranges
------

A `Range` is a set of conditions which specify which versions satisfy the range.

A condition is composed of an operator and a version. The supported operators are:

- `<1.0.0` Less than `1.0.0`
- `<=1.0.0` Less than or equal to `1.0.0`
- `>1.0.0` Greater than `1.0.0`
- `>=1.0.0` Greater than or equal to `1.0.0`
- `1.0.0`, `=1.0.0`, `==1.0.0` Equal to `1.0.0`
- `!1.0.0`, `!=1.0.0` Not equal to `1.0.0`. Excludes version `1.0.0`.

Note that spaces between the operator and the version will be gracefully tolerated.

A `Range` can link multiple `Ranges` separated by space:

Ranges can be linked by logical AND:

  - `>1.0.0 <2.0.0` would match between both ranges, so `1.1.1` and `1.8.7` but not `1.0.0` or `2.0.0`
  - `>1.0.0 <3.0.0 !2.0.3-beta.2` would match every version between `1.0.0` and `3.0.0` except `2.0.3-beta.2`

Ranges can also be linked by logical OR:

  - `<2.0.0 || >=3.0.0` would match `1.x.x` and `3.x.x` but not `2.x.x`

AND has a higher precedence than OR. It's not possible to use brackets.

Ranges can be combined by both AND and OR

  - `>1.0.0 <2.0.0 || >3.0.0 !4.2.1` would match `1.2.3`, `1.9.9`, `3.1.1`, but not `4.2.1`, `2.1.1`

Note that all ranges can contain a mixture of 3 and 4 digit version numbers.

Range usage:

```
v, err := semver.Parse("1.2.3")
expectedRange, err := semver.ParseRange(">1.0.0 <2.0.0 || >=3.0.0")
if expectedRange(v) {
    //valid
}

```

Example
-----

Have a look at full examples in [v4/examples/main.go](v4/examples/main.go)

```go
import github.com/khulnasoft/semver/v4

v, err := semver.Make("0.0.1-alpha.preview+123.github")
fmt.Printf("Major: %d\n", v.Major)
fmt.Printf("Minor: %d\n", v.Minor)
fmt.Printf("Patch: %d\n", v.Patch)
fmt.Printf("Pre: %s\n", v.Pre)
fmt.Printf("Build: %s\n", v.Build)

// Prerelease versions array
if len(v.Pre) > 0 {
    fmt.Println("Prerelease versions:")
    for i, pre := range v.Pre {
        fmt.Printf("%d: %q\n", i, pre)
    }
}

// Build meta data array
if len(v.Build) > 0 {
    fmt.Println("Build meta data:")
    for i, build := range v.Build {
        fmt.Printf("%d: %q\n", i, build)
    }
}

v001, err := semver.Make("0.0.1")
// Compare using helpers: v.GT(v2), v.LT, v.GTE, v.LTE
v001.GT(v) == true
v.LT(v001) == true
v.GTE(v) == true
v.LTE(v) == true

// Or use v.Compare(v2) for comparisons (-1, 0, 1):
v001.Compare(v) == 1
v.Compare(v001) == -1
v.Compare(v) == 0

// Manipulate Version in place:
v.Pre[0], err = semver.NewPRVersion("beta")
if err != nil {
    fmt.Printf("Error parsing pre release version: %q", err)
}

fmt.Println("\nValidate versions:")
v.Build[0] = "?"

err = v.Validate()
if err != nil {
    fmt.Printf("Validation failed: %s\n", err)
}
```


Benchmarks
-----

    BenchmarkParseSimple-8             1604840    757.7 ns/op    224 B/op   2 allocs/op
    BenchmarkParseComplex-8             480861     2351 ns/op    434 B/op   8 allocs/op
    BenchmarkParseAverage-8             568860     2365 ns/op    341 B/op   5 allocs/op
    BenchmarkParseTolerantAverage-8     423841     2881 ns/op    535 B/op   9 allocs/op
    BenchmarkStringSimple-8           38751848    31.04 ns/op      5 B/op   1 allocs/op
    BenchmarkStringLarger-8           16473154    71.87 ns/op     32 B/op   2 allocs/op
    BenchmarkStringComplex-8          10292580    153.5 ns/op     80 B/op   3 allocs/op
    BenchmarkStringAverage-8          11083524    109.2 ns/op     48 B/op   2 allocs/op
    BenchmarkValidateSimple-8        300665788    3.944 ns/op      0 B/op   0 allocs/op
    BenchmarkValidateComplex-8         6201321    217.8 ns/op      0 B/op   0 allocs/op
    BenchmarkValidateAverage-8        11562688    106.9 ns/op      0 B/op   0 allocs/op
    BenchmarkCompareSimple-8         175868226    6.688 ns/op      0 B/op   0 allocs/op
    BenchmarkCompareComplex-8         81058132    15.18 ns/op      0 B/op   0 allocs/op
    BenchmarkCompareAverage-8         51633844    23.23 ns/op      0 B/op   0 allocs/op
    BenchmarkSort-8                    6390456    190.1 ns/op    264 B/op   2 allocs/op
    BenchmarkRangeParseSimple-8        1061410     1169 ns/op    418 B/op   8 allocs/op
    BenchmarkRangeParseAverage-8        504180     2263 ns/op    844 B/op  15 allocs/op
    BenchmarkRangeParseComplex-8        170269     7123 ns/op   2906 B/op  45 allocs/op
    BenchmarkRangeMatchSimple-8       71724456    16.70 ns/op      0 B/op   0 allocs/op
    BenchmarkRangeMatchAverage-8      32207384    36.53 ns/op      0 B/op   0 allocs/op
    BenchmarkRangeMatchComplex-8      12062602    100.4 ns/op      0 B/op   0 allocs/op

See benchmark cases at [semver_test.go](semver_test.go)


Motivation
-----

I simply couldn't find any lib supporting the full spec. Others were just wrong or used reflection and regex which i don't like.

Motivation for fork
-----

Khulnasoft Gateway and Khulnasoft Gateway OSS use different versioning systems; Khulnasoft Gateway doesn't follow strict semver. As Khulnasoft Gateway and
Khulnasoft Gateway OSS are released the major.minor.patch versions are in lock step; however, Khulnasoft Gateway may have additional revisions
that are based of the OSS upstream. In order to handle version comparisons the fork was a necessary evil.

Contribution
-----

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.


License
-----

See [LICENSE](LICENSE) file.
