export const getJWT = () => {
    const item = localStorage.getItem("jwt")
    return JSON.parse(item)
}

export const getUsername = () => {
    return getJWT()['Username']
}

export const getRoles = () => {
    return getJWT()['Roles']
}

export const hasRole = (name) => {
    return getRoles() || [].indexOf(name) >= 0
}