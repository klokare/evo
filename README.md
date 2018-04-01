evo 
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.png?v=103)](https://opensource.org/licenses/mit-license.php) [![Build Status](https://travis-ci.org/klokare/evo.svg?branch=master)](https://travis-ci.org/klokare/evo) [![Coverage Status](https://coveralls.io/repos/github/klokare/evo/badge.svg)](https://coveralls.io/github/klokare/evo) [![Go Report Card](https://goreportcard.com/badge/github.com/klokare/evo)](https://goreportcard.com/report/github.com/klokare/evo) [![GoDoc](https://godoc.org/github.com/klokare/evo?status.svg)](https://godoc.org/github.com/klokare/evo)
==============

**evo** is neuroevolution framework based upon Dr. Kenneth Stanley's [NEAT](https://www.cs.ucf.edu/~kstanley/neat.html) and subsequent extensions. Built from the ground up from the research papers and online articles, this implementation strives to be performant and extensible. See [the wiki](https://github.com/klokare/evo/wiki) for more details.

> NOTE: This is the second incarnation of EVO. Having survived the growing pains of the original, I decided to update the library based on my experience and continued reading. The prior version is archived under the archive-20180109 branch.

## Installing
To start using EVO, install Go and run `go get`:

```sh
$ go get github.com/klokare/evo/...
```

For further information on using, see the examples and peruse [the wiki](https://github.com/klokare/evo/wiki).

## Version history and upcoming releases
Version|Description
-------|-----------
0.1|core library and tests (completed)
0.2|default configurer (completed)
0.3|default network and translator (completed)
0.4|NEAT-equivalent package and XOR experiment (completed)
0.5|~~phased mutator and OCR experiment (completed)~~ _temporarily removed until better OCR experiment implemented__
0.6|HyperNEAT package and boxes experiment (completed)
0.7|*configurer rewrite, removal of functional options*
0.8|ES-HyperNEAT package and mazes experiment
0.9|novelty package and updated mazes experiment
0.10|real-time package
1.0|production-ready release