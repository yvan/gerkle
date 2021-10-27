gerkle
--

gerkle is a merkle tree implementation in go. i decided to use go because i'm looking for substitutes for python (which i still want to use) that are a bit more reliable and don't just typecast all over the place, as if types were cranberry sauce at your aunt's thanksgiving dinner.

why didn't i use cmp for the tests?

well i felt writing my own function to compare the structs was easier than debugging weird shit like "panic: cannot handle unexported field at {main.Node}.hash:" or issues with go modules, GOPATH, etc.chock it up to beginners naivety.

micro todo:

1 - function to create lists of leaves  (f)
2 - function to create trees (f)
3 - function to pretty print the trees 
4 - we need a function that can check trees to verify if they've changed and quickly locate the changed leaves
5 - 