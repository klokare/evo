Boxes Visual Discrimination
===========================

This example comes from the visual discrimination task described in [A Hypercube-Based Indirect Encoding for Evolving Large-Scale Neural Networks](http://eplex.cs.ucf.edu/papers/stanley_alife09.pdf) (Stanley, et. al.). The paper utilises this task to demonstate HyperNEAT's ability to incorporate the problem's geometric information. For EVO, this tasks represents an ideal problem to demonstrate the differences in NEAT, HyperNEAT, HyperNEAT-LEO, ES-HyperNEAT, and ES-HyperNEAT-LEO. 

## Visualisng the best solution
When performing a single run of the experiment, an image file is created (use the `--image` flag to specify the output filename) showing, for each test case, the locations of the small and large boxes as well as the network's guess at the centre of the large box. 

## TODO
• add HyperNEAT experiment
• add ES-HyperNEAT experiment
• add archive and restore helpers
• use restore to enable running previous best at different resolutions