{
  mkIndex(arr)::
    local idxs = [ "[%s]" % i for i in arr];
    "." +std.join(".", idxs),

  mkKey(k):: if k == "." then "" else k,

  indexify(obj, index):: {
    [$.mkIndex(index) + $.mkKey(key)]: obj[key],
    for key in std.objectFields(obj)
  },
}
