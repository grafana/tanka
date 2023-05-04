# Explanation

This tests the weird case where relative imports can be resolved either from the place where they are defined or from the main.jsonnet being run.
The code has to consider all files within the environment context as being potentially imported and used by itself

*Note that this is not a good practice, but Tanka's behavior should be consistent with the Jsonnet interpreter*
