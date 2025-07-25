
function >> test -> int {
    int x;
    0 -> x;
}

function int x >> f -> float {
    int y;
    x -> y;
    y + 1 -> y;
}

function string s, int x >> g -> bool {
    0 => i;
    while i < x {
        s >> print;
        i + 1 -> i;
    }
}




// // Want:
// function ((int, int -> bool) f, int x, int y) >> h -> int {
//     if (x, y >> f){
//         "f(x, y) is positive" >> print;
//     } else {
//         "f(x, y) is negative" >> print;
//     }
// }

// function () >> f -> int {
//     1 >> return;
// }