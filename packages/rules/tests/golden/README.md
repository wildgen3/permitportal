# Golden tests

Profile in, expected obligation set out. **This is the correctness specification, not a
regression suite.** A change that alters a golden expectation requires the pull request
to say why the law changed — or to admit the rule was wrong.

`false_negatives` is the gated metric: an obligation that should apply and does not
appear is a release blocker. Extra obligations are a defect too, but a different and less
harmful one.
