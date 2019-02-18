# EAI

EAI is an acronym for EAI Ain't Interest, or something even more forgettable<sup><a name="eai_footnote_link" href="#eai_footnote_body">1</a></sup>. It is an incentive for holding lots of ndau, plus an incentive for locking ndau. These incentives take the form of ndau, conjured from the ether. They are awarded when an account's delegated node decides to issue a `CreditEAI` transaction, on whatever schedule that node feels like.

## But how is it calculated?

Use the [`eai.Calculate`](https://github.com/oneiro-ndev/ndaumath/blob/cf6f1e6fc1f3a54925c7f82a670cbeb11ae49ebe/pkg/eai/eai.go#L11-L43) function.

## No really, how do I calculate it by hand?

If you absolutely must hand-calculate EAI for verification or other purposes, take as your first reference the [test cases](https://github.com/oneiro-ndev/ndaumath/blob/cf6f1e6fc1f3a54925c7f82a670cbeb11ae49ebe/pkg/eai/eai_test.go#L241-L449), which are fairly well documented.

For any account, over any period of time, one can break down the period into some number of `(rate, duration)` pairs. Once you have the correct list of pairs, you simply compute the product <a href="https://www.codecogs.com/eqnedit.php?latex=\mathrm{factor}&space;=&space;\prod{e^{\mathrm{rate}&space;*&space;\mathrm{duration}}}" target="_blank"><img src="https://latex.codecogs.com/gif.latex?\mathrm{factor}&space;=&space;\prod{e^{\mathrm{rate}&space;*&space;\mathrm{duration}}}" title="\mathrm{factor} = \prod{e^{\mathrm{rate} * \mathrm{duration}}}" /></a>, where `rate` is expressed numerically (i.e. 1% == `0.01`) and `duration` is expressed as a fractional year (i.e. 1 day == `1/365`), using the [ndau canonical duration definitions](https://github.com/oneiro-ndev/ndaumath/blob/cf6f1e6fc1f3a54925c7f82a670cbeb11ae49ebe/pkg/types/duration.go#L175-L195). After computing `factor`, <a href="https://www.codecogs.com/eqnedit.php?latex=\mathrm{eai}&space;=&space;\mathrm{initial\_balance}&space;*&space;(\mathrm{factor}&space;-&space;1)" target="_blank"><img src="https://latex.codecogs.com/gif.latex?\mathrm{eai}&space;=&space;\mathrm{initial\_balance}&space;*&space;(\mathrm{factor}&space;-&space;1)" title="\mathrm{eai} = \mathrm{initial\_balance} * (\mathrm{factor} - 1)" /></a>.

We've glossed over the mechanism for getting the `(rate, duration)` pairs, because it's complicated.

### Computing `(rate, duration)` pairs for an arbitrary period

The easiest portion of EAI rate to calculate has to do with the lock: if an account is locked, then at the time of lock, a bonus lock rate is computed by reference to a lock rate lookup table, and stored with the lock. For a locked account, simply retrieve the bonus lock rate. This will be added to all other rates computed.

Next, compute the effective age of the account for the period in question. The effective age of the account is the weighted average age plus the lock period.

Next, decompose the period into some number of periods based on the unlocked rate table. Note: if the account is notified, the effective age freezes at that point and does not change until the account is unlocked.

A worked example is likely to be helpful here. This is from [one of the test cases](https://github.com/oneiro-ndev/ndaumath/blob/cf6f1e6fc1f3a54925c7f82a670cbeb11ae49ebe/pkg/eai/eai_test.go#L330-L371):

```
Case 3: What happens if an account is:

- locked for 180 days
- notified to unlock 165 days from now
- 84 days since last EAI update
- current actual weighted average age is 123 days

The span of effective average age we care about for the unlocked
portion runs from actual day 39 to actual day 123. The notify happens
on actual day 108. It expires on actual day 288. At that point, the
rate will drop back to the actual weighted average age.

The effective period begins on day 129, and runs forward normally
until effective day 157. Effective time freezes at that point. On
actual day 157, the notice period ends and calculations resume using
the actual weighted average age.

Dashed lines in the following graph indicate points in the future,
assuming no further transactions are issued.

10%                     ┌────────|────────x-------
 9%              ┌──────┘        |
 8%      ──x─────┘               |
        ___________________________________________
 actual    39    60     90      108      123   288
 effect.  219   240    270      288......288...288
 month    (7)   (8)    (9)                     (9)

Because the account was locked for 180 days, and 180 days has a bonus
rate of 2%, the actual rate used during the lock and notification
periods should increase by a constant rate of 2%.
We thus get the following calculation to compute the EAI factor:

   e^(10% * 21 days)
 * e^(11% * 30 days)
 * e^(12% * 33 days)

The 33 days of the final term are simply the 18 unnotified days
of the rate period plus the 15 days notified to date.
```

---
<small>
<a name="eai_footnote_body" href="#eai_footnote_link">1</a>: Ecosystem Alignment Incentive
</small>
