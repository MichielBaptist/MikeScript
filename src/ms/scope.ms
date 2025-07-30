
"Hello" => s;

{

    function () >> s_printer -> (->) {
        function () >> g {
            s >>= print;
        }
        return g;
    }

    // Empty function variable
    var (->) my_first_s_printer;

    // Print hello
    "hello" >>= print;
    =s_printer -> my_first_s_printer;
    "hello" >>= print;
    if true {
        =s_printer -> my_first_s_printer;
    }
    
    // Call my printer
    =my_first_s_printer;

    // Change s to world
    "World" => s;

    =my_first_s_printer; // expected "Hello"

    // Now s should be in this scope
    =s_printer => my_second_s_printer;

    // Call my second printer
    =my_second_s_printer;

    =print_env;
}
