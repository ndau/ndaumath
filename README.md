# ndaumath
Definitive math libraries for calculations related to ndau


## Purpose

This is intended to be the canonical location for all of the mathematical calculations for ndau.

The reason for this is that we are not just adding and subtracting -- we have EAI and SIB calculations that involve fractions, and we have to be very careful about things like overflow and precision.

We also use this library to define a robust test suite that is intended to be independent of implementation language. Thus, the same test suite can be used to prove that the math is correct no matter how it is implemented.

Finally, this repository should contain the ultimate reference documentation for ndau's mathematics.

