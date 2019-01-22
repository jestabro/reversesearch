# reversesearch

The package reversearch provides an interface for exploring the
Reverse Search algorithm of Avis and Fukuda:

https://www.sciencedirect.com/science/article/pii/0166218X9500026N

The preliminary implementation here is drawn from David Avis'
tutorial:

http://cgm.cs.mcgill.ca/~avis/doc/tutorial1/tutorial.pdf

The algorithm was originally invented to investigate convex hull
calculations, but is amenable to calculation in a variety of 
combinatorial settings.

The first example here enumerates all permutations, given length n; 
others to follow.

The purpose of this implementation is to take advantage of Go's
flexible interface structure to untangle the original C
implementation, find the proper level of generalization, and refactor
to take advantage of the inherent recursive and concurrent nature of 
the algorithm.
