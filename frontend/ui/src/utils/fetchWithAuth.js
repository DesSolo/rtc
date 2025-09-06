export const fetchWithAuth = (url, options = {}) => {
    const token = localStorage.getItem('token');

    const headers = {
        ...options.headers,
        'Authorization': `jwt ${token}`
    };

    return fetch(url, { ...options, headers });
};