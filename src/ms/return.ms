
function () >> f {
    "Hello " + =g >>= print;
    return;
    "Goodbye" >>= print;
}

function () >> g -> string {
    return "world";
}

// call
=f;
=g >>= print;