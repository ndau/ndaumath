address = require("./keyaddr.js");

mk = address.NewPrivateMaster('asdfwdsdfdewidkdkffjfjfuggujfjug');
console.log(mk);
mk.then((k) => {
    console.log(k)
    pubkey = k.Neuter()
    console.log(pubkey);
}).catch((e) => console.log("error: ", e));

mk.then((k) => {
    msg = "CAFEF00DBAAD1DEA"
    k.Sign(msg).then((sig) => {
        console.log(sig)
        pubkey = k.Neuter()
        sig.Verify(msg, k.key).then((b) => {
            console.log("Verified: ", b)
        }).catch((e) => console.log("No verify!", e));
    }).catch((e) => console.log("sigerror: ", e));

}).catch((e) => console.log("error: ", e));
