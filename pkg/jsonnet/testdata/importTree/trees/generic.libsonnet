{
  new(breed, color, size='m'):: {
    kind: 'tree',

    breed: breed,
    color: color,
    size: size,

    needs: 'water',
    eats: 'co2',
    creates: 'o2',
    keeps: 'the world healthy',
    nested_random_struct:: {
      branch: {
        branch: {
          branch1: {
            branch: import 'leaf.libsonnet',
          },
          branch2: {
            branch: {
              branch1: {
                branch: {
                  branch: import 'leaf.libsonnet',
                },
              },
              branch2: {
                branch: {
                  branch1: {
                    branch: import 'leaf.libsonnet',
                  },
                  branch2: {
                    branch: {
                      branch: {
                        branch: {
                          branch: import 'leaf.libsonnet',
                        },
                      },
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
}
