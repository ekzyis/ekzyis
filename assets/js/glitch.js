// this file creates a glitch effect by switching randomly between two titles

const title1 = "Javascript Fatigue"
const title2 = "HTMX Intrigue"

function randInt(n) {
    return Math.floor(Math.random() * n)
}

async function glitch(node) {
    if (!node) return

    const r1 = randInt(10)
    for (let i = 0; i < r1; ++i) {
        const r2 = randInt(100)
        await new Promise(resolve => {
            setTimeout(() => {
                node.innerHTML = node.innerHTML === title1 ? title2 : title1;
                resolve()
            }, r2)
        })
    }
}

const titleNode = document.querySelector("h1.post-title")
const listNode = document.querySelector("h4.post-item-title > a")

glitch(titleNode)
glitch(listNode)