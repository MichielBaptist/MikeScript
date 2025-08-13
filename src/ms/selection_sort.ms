
function ([]float a, int l) >> selection_sort -> []float {
    if l <= 1 {
        return a;
    }
    
    0 => i;
    while i < l {

        i => min_idx;
        a[min_idx] => min;

        // look for smallest
        i + 1 => j;
        while j < l {
            if a[j] < min {
                a[j -> min_idx] -> min;
            }
            j + 1 -> j;
        }

        // swap
        a[i] => tmp;
        a[min_idx] -> a[i];
        tmp -> a[min_idx];

        i + 1 -> i;
    }

    return a;
}

function (int len) >> rand_array -> []float {

    // allocate floats
    [len]float{} => a;

    0 => i;
    while i < len {
        =rand -> a[i];
        i + 1 -> i;
    }
    return a;
}

10000 => n;
(n >>= rand_array), n >>= selection_sort >>= print;