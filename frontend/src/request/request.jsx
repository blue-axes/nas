

const options = {
    credentials: "same-origin",
    mode: "cors",
    redirect: "follow",
    referrer: "",
}

function Download(url, requestOptions, callback) {
    fetch(url, { ...options, ...requestOptions, method: "GET" }).
        then(resp => resp.blob()).
        then(blob => {
            const link = document.createElement('a')
            link.style.display = "none"
            link.href = URL.createObjectURL(blob)
            document.body.appendChild(link)
            link.click()
            URL.removeObjectURL(link.href)
            document.body.removeChild(link)
            callback('')
        }).catch(err => {
            callback(err)
        })
}

function Upload(url, requestOptions) {
    
}
