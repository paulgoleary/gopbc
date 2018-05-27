# gopbc

This is a POC implementation of pairing-based cryptography which is inspired by Ben Lynn's **pbc** project: https://crypto.stanford.edu/pbc/.
For ideas / testing / inspiration I also studied JPBC (http://gas.dia.unisa.it/projects/jpbc/#.WwrPhVMvw3E) - a Java implementation of pbc - and the amazing RELIC toolkit (https://github.com/relic-toolkit), which implements just about everything.

The current work only implements one specific pairing type: what pbc refers to as a type A curve with the Tate pairing.
Some of the more standard optimizations are also present such as evaluating the Miller loop in projective coordinates with NAF - nothing novel there.

As I also note in a separate POC project (https://github.com/paulgoleary/goUmbral), to the best of my knowledge there does not seem to be a general Golang implementation of the Weierstrass curve.
This caused me to borrow my implementation here for that other project. I likely should factor it out into it's own shared project as it may be useful on it's own.
The base curve here is actually fairly well tested for correctness and compatibility with other implementations, eg. pbc and RELIC.

For validation purposes the pairing is used to demonstrate the basics of AFGH-style proxy-reencryption and BLS signatures.

I have paused on this work. In other research, I came across the excellent crypto work being done by the folks at Zcash.
Crypto research is evolving and progressing very rapidly.
It's worth noting that the curve types currently implemented in pbc (and JPBC) may now be considered suboptimal WRT their security level and performance.
This observation has caused Zcash to begin to re-implement with a different, stronger class of curve for their pairing based crypto.
Through extensive research they settled on and implemented a curve type called BLS12-381: https://blog.z.cash/new-snark-curve/.

In this project I began to research a parallel implementation of BLS12-381 in Golang. The main difference with BLS12-381 is the higher-order (degree 12) extension field.
In BLS12-381 this is constructed with towers over polynomial fields.

The world of crypto benefits greatly from standardization. It appears that the Zcash team coordinated with the author of RELIC and that tookit also now implements BLS12-381.
It may be useful to have a pure Golang implementation as well, as proposed here. If interested, please reach out. In the meantime I decided a better use of my time would be to learn Rust. :)
