var array = []

var intervals = {
    add: async (i) => {
        array.push(i)
    },

    clear: async (i) => {
        clearInterval(i)
        array = array.filter(item => item != i)
    },

    clearAll: async () => {
        array.forEach((i) => {
            clearInterval(i)
        })

        array = []
    }
}

export default intervals