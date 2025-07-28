
// function with no arguments and no return type
function () >> no_args_no_return {
    
}

// function with no arguments
function () >> no_args_with_return -> int {
    return 0;
}

// Function with one argument and no return type
function (int x) >> one_arg_no_return {
    x + 10 >>= print;
}

// Function with one argument and a return type
function (int x) >> one_arg_with_return -> int {
    return x + 1;
}

function (int y, int x) >> is_multiple -> bool {
    return x % y == 0;
}
function (int x) >> is_even -> bool {
    return 2, x >>= is_multiple;
}

// Function with function as an argument
function ((int -> bool) f, int x, int y) >> first_order_function {
    "f(x) is", (x >>= f) >>= print;
    "f(y) is", (y >>= f) >>= print;
}

// Function returning a function
function (int y) >> is_multiple_of -> (int -> bool) {
    return y >> is_multiple;
}

// Call simple
"expect:" >>= print;
"nothing" >>= print;
"got: " >>= print;
=no_args_no_return >>= print;
"expect:" >>= print;
"0" >>= print;
"got: " >>= print;
=no_args_with_return >>= print;
"expect:" >>= print;
"15" >>= print;
"got: " >>= print;
5 >>= one_arg_no_return;
"expect:" >>= print;
"6" >>= print;
"got: " >>= print;
5 >>= one_arg_with_return >>= print;
"expect: " >>= print;
"f(x) is, true" >>= print;
"f(y) is, false" >>= print;
"got: " >>= print;
is_even, 4, 5 >>= first_order_function;
"expect: " >>= print;
"f(x) is, true" >>= print;
"f(y) is, false" >>= print;
"got: " >>= print;
(3>>=is_multiple_of), 6, 7 >>= first_order_function;