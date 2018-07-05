address = require("./keyaddr.js");

mk = address.NewPrivateMaster('asdfwdsdfdewidkdkffjfjfuggujfjug');
console.log(mk);
mk.then((k) => {
    console.log(k)
    pubkey = k.Neuter()
    console.log(pubkey);
}).catch((e) => {
    console.log("master key error: ", e);
    process.exit(1);
});

mk.then((k) => {
    msg = "CAFEF00DBAAD1DEA"
    k.Sign(msg).then((sig) => {
        console.log(sig)
        pubkey = k.Neuter()
        sig.Verify(msg, k.key).then((b) => {
            console.log("Verified: ", b)
        }).catch((e) => {
            console.log("NO verify: ", e);
            process.exit(1);
        });

    }).catch((e) => {
        console.log("sigerror: ", e);
        process.exit(1);
    });

}).catch((e) => {
    console.log("error: ", e);
    process.exit(1);
});
