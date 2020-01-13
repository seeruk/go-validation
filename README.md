# go-validation

Simple and flexible Go (Golang) validation library.

## TODO

### Finish constraints:

* AnyNRequired
* Equals
* ExactlyNRequired
* ~~Length~~
* Max
* ~~MaxLength~~
* Min
* ~~MinLength~~
* ~~MutuallyExclusive~~
* MutuallyInclusive
* Nil
* NotEquals
* NotNil
* OneOf
* Predicate
* Regexp
* RegexpString
* ~~Required~~
* TimeAfter
* TimeBefore
* TimeStringAfter
* TimeStringBefore
* Valid

### Tests

* Write 'Max' tests.
* Tests should also try pointers to types with nil to verify they don't panic when given nil, many
constraints should just be optional in that case.
