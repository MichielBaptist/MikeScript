

// function (int x) >> outer -> int {

//     // declare inner function
//     function (int y) >> inner -> int {

//         // declare inner inner function
//         function (int z) >> inner_inner -> int {
//             return z + 1;
//         }

//         return y + 1 >>= inner_inner;
//     }

//     return x + 1 >>= inner;
// }

function (int x) >> x_printer -> (->) {

    -1 => i;

    // x printer
    function () >> f {
        "Iteration", (i + 1 -> i) >>= print;
        (x * 2) + 1 -> x >>= print;
    }

    return f;
}


// expected to be 4
// 1 >>= outer >>= print;

// get an x printer (will return error normally saying x is not defined)
0 >>= x_printer => my_printer;

-1 => i;
while (i + 1 -> i) < 10 {
    =my_printer;
}