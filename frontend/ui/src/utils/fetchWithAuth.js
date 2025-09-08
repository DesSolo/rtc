export const fetchWithAuth = async (url, options = {}, navigate) => {
    const token = localStorage.getItem('token');

    const headers = {
        ...options.headers,
        'Authorization': `jwt ${token}`
    };

    const response = await fetch(url, { ...options, headers });

    if (response.status === 401) {
        if (navigate) {
            navigate('/login');
        } else {
            window.location.href = '/login';
        }
    }

    return response;
};