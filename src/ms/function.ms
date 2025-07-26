
// function with no arguments and no return type
function () >> no_args_no_return {}

// function with no arguments
function () >> no_args_with_return -> int {
    0;
}

// Function with one argument and no return type
function (int x) >> one_arg_no_return {
    x >>= print;
}

// Function with one argument and a return type
function (int x) >> one_arg_with_return -> int {
    x + 1;
}

function (int x) >> is_even -> bool {
    2, x >>= is_multiple;
}

// Function with function as an argument
function ((int -> bool) f, int x, int y) >> first_order_function {
    "f(x) is", (x >>= f) >>= print;
    "f(y) is", (y >>= f) >>= print;
}

function (int y, int x) >> is_multiple -> bool {
    x % y == 0;
}

// Function returning a function
function (int y) >> is_multiple_of -> (int -> bool) {
    y >> is_multiple;
}
