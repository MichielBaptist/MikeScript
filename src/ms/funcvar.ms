var (->) hello;

function () >> f -> (->) {
    function () >> g {
        "Hello world" >>= print;
    }
    return g;
}

=f -> hello;

// print hello world
=hello;