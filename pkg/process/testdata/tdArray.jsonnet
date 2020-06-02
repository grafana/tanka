local utils = import './utils.libsonnet';

local deep = import './tdDeep.jsonnet';
local flat = import './tdFlat.jsonnet';

{
  deep: [deep.deep, [flat.deep]],
  flat: utils.indexify(deep.flat, [0]) +
        utils.indexify(flat.flat, [1, 0]),
}
