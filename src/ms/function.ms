
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
    string t;
}