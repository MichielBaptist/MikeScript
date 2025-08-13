


function ([]float a, int start, int end) >> quick_sort -> []float {

    if start >= end {
        return a;
    }

    a, start, end >>= partition => pivot;

    a, start, pivot - 1 >>= quick_sort;
    a, pivot + 1, end >>= quick_sort;

    return a;

}

function ([]float a, int start, int end) >> partition -> int {

    // select pivot
    a[end] => pivot;

    start => i;
    start => j;

    while j < end {

        if a[j] < pivot {
            a, i, j >>= swap;
            i + 1 -> i;
        }

        j + 1 -> j;
    }

    a, i, end >>= swap;

    return i;
}

function ([]float a, int i, int j) >> swap {
    a[i] => tmp;
    a[j] -> a[i];
    tmp -> a[j];
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

100000 => n;
(n >>= rand_array), 0, n-1 >>= quick_sort >>= print;