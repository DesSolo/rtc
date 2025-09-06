export const getJWT = () => {
    const item =localStorage.getItem("jwt")
    return JSON.parse(item)
}

export const getUsername = () => {
    return getJWT()['Username']
}