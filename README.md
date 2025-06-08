# LINC

## LINC Is Not a Compiler

LINC is a program that transpiles Scratch projects to native Go code.

## What about `sb-edit` (`Leopard`'s Transpiler)

I do know about `sb-edit` and `Leopard` but I also understand that they have their downsides. The main one is that even though `Leopard` is theoretically faster than Scratch, it is still running as Javascript in a browser. This means that `Leopard` is still pretty slow especially for larger and more complicated projects like 3D Renderers and AI. The goal of LINC is by transpiling to a language that can be natively compiled, to run Scratch projects at speeds faster than previously possible (at least for Scratch 3.0.) There are some downsides to this; one of which is that there are two steps to compile with LINC, first you have to transpile to Go, then you have to compile the transpiled Go code into a native executable. Both projects have their uses and LINC is not mean to replace `sb-edit` or `Leopard` in any way.
