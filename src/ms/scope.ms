
"Hello" => s;

{

    function () >> prints {
        =print_env;
        s >>= print;
    }

    // Function is expected to capture closure statically
    =prints; // expected "Hello"
    "World" => s;
    =prints; // expected "Hello"
}