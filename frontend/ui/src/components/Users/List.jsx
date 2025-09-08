import { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { fetchWithAuth } from "../../utils/fetchWithAuth.js";
import { Button, Flex, Input, Modal, Switch, Table, message, Tag, Tooltip, Space } from "antd";
import { PlusOutlined, UserOutlined } from "@ant-design/icons";
import { useOutletContext, useSearchParams } from "react-router-dom";
import CreateUser from "./Create.jsx";
import { getUsername } from "../../utils/storage.js";

const UsersList = () => {
    const { setTitle } = useOutletContext();
    const navigate = useNavigate();
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(false);
    const [totalUsers, setTotalUsers] = useState(0);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [searchParams, setSearchParams] = useSearchParams();

    // Get initial parameters from URL
    const paramPage = parseInt(searchParams.get("page") || "1", 10);
    const paramLimit = parseInt(searchParams.get("limit") || "10", 10);
    const paramQ = searchParams.get("q") || "";

    const [searchValue, setSearchValue] = useState(paramQ);
    const [currentPage, setCurrentPage] = useState(paramPage);
    const [pageSize, setPageSize] = useState(paramLimit);
    const searchTimeoutRef = useRef(null);

    // Компонент для отображения ролей
    const RolesDisplay = ({ roles }) => {
        if (!roles || !Array.isArray(roles) || roles.length === 0) {
            return <Tag icon={<UserOutlined />} color="default">No roles</Tag>;
        }

        const visibleRoles = roles.slice(0, 3);
        const hiddenRolesCount = roles.length - 3;

        return (
            <Space size={[0, 4]} wrap>
                {visibleRoles.map((role, index) => {
                    let color = 'blue';

                    if (role.toLowerCase().includes('admin')) {
                        color = 'red';
                    } else if (role.toLowerCase().includes('user')) {
                        color = 'green';
                    } else if (role.toLowerCase().includes('editor')) {
                        color = 'orange';
                    } else if (role.toLowerCase().includes('viewer')) {
                        color = 'purple';
                    }

                    return (
                        <Tooltip key={index} title={role}>
                            <Tag
                                color={color}
                                icon={<UserOutlined />}
                                style={{
                                    margin: 0,
                                    maxWidth: 100,
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                    cursor: 'pointer'
                                }}
                            >
                                {role}
                            </Tag>
                        </Tooltip>
                    );
                })}

                {hiddenRolesCount > 0 && (
                    <Tooltip title={roles.slice(3).join(', ')}>
                        <Tag style={{ cursor: 'pointer' }}>
                            +{hiddenRolesCount}
                        </Tag>
                    </Tooltip>
                )}
            </Space>
        );
    };

    const columns = [
        {
            title: "Username",
            dataIndex: "username",
            key: "username"
        },
        {
            title: "Enabled",
            dataIndex: "is_enabled",
            key: "is_enabled",
            render: (is_enabled, record) => (
                <Switch
                    disabled={record.username === getUsername()}
                    checked={is_enabled}
                    onChange={(checked) => handleEnableChange(record.username, checked)}
                />
            )
        },
        {
            title: "Roles",
            dataIndex: "roles",
            key: "roles",
            render: (roles) => <RolesDisplay roles={roles} />
        },
        {
            title: "Created",
            dataIndex: "created_at",
            key: "created",
            render: (text) => <>{new Date(text).toLocaleString()}</>,
        }
    ];

    const fetchUsers = (page = currentPage, limit = pageSize, q = searchValue) => {
        setLoading(true);

        const offset = (page - 1) * limit;
        const params = new URLSearchParams({
            limit: limit.toString(),
            offset: offset.toString(),
        });

        if (q) {
            params.append('q', q);
        }

        fetchWithAuth(`/api/v1/users?${params}`, {}, navigate)
            .then((response) => {
                if (!response.ok) throw new Error("Ошибка загрузки данных");
                return response.json();
            })
            .then((data) => {
                setUsers(data.data.users);
                setTotalUsers(data.data.total);

                // Обновляем параметры URL
                const newSearchParams = new URLSearchParams();
                newSearchParams.set('page', page.toString());
                newSearchParams.set('limit', limit.toString());
                if (q) {
                    newSearchParams.set('q', q);
                }
                setSearchParams(newSearchParams);
            })
            .catch((error) => {
                console.error('Error fetching users:', error);
                message.error('Ошибка при загрузке пользователей');
            })
            .finally(() => {
                setLoading(false);
            });
    };

    const handleTableChange = (pagination, filters, sorter) => {
        const { current, pageSize } = pagination;
        setCurrentPage(current);
        setPageSize(pageSize);
        fetchUsers(current, pageSize, searchValue);
    };

    const handleSearch = (value) => {
        setSearchValue(value);
        setCurrentPage(1);
        fetchUsers(1, pageSize, value);
    };

    const handleSearchChange = (e) => {
        const value = e.target.value;
        setSearchValue(value);

        // Очищаем предыдущий таймаут
        if (searchTimeoutRef.current) {
            clearTimeout(searchTimeoutRef.current);
        }

        // Устанавливаем новый таймаут для поиска через 500 мс
        searchTimeoutRef.current = setTimeout(() => {
            handleSearch(value);
        }, 500);
    };

    const handleEnableChange = (username, enabled) => {
        fetchWithAuth(`/api/v1/users/${username}`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ is_enabled: enabled }),
        }, navigate)
        .then((response) => {
            if (!response.ok) throw Error("error")

            fetchUsers(paramPage, paramLimit, paramQ)
        })
    };

    useEffect(() => {
        fetchUsers(paramPage, paramLimit, paramQ);

        // Очищаем таймаут при размонтировании компонента
        return () => {
            if (searchTimeoutRef.current) {
                clearTimeout(searchTimeoutRef.current);
            }
        };
    }, []);

    useEffect(() => {
        setTitle("Users");
    }, [setTitle]);

    return (
        <>
            <Flex justify="space-between" style={{ marginBottom: 16 }}>
                <Input
                    placeholder="Search by username"
                    value={searchValue}
                    onChange={handleSearchChange}
                    style={{ marginRight: 16, width: 200 }}
                    allowClear
                />
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => setIsModalOpen(true)}
                >
                    New User
                </Button>
            </Flex>

            <Table
                columns={columns}
                dataSource={users}
                rowKey="username"
                loading={loading}
                pagination={{
                    current: currentPage,
                    pageSize: pageSize,
                    total: totalUsers,
                    showSizeChanger: true,
                    pageSizeOptions: ['10', '20', '50'],
                    showTotal: (total, range) =>
                        `${range[0]}-${range[1]} of ${total} users`,
                }}
                onChange={handleTableChange}
            />

            <Modal
                title="Create new user"
                open={isModalOpen}
                onCancel={() => setIsModalOpen(false)}
                width={600}
                footer={[
                    <Button form="createUser" type="primary" key="submit" htmlType="submit">
                        Create
                    </Button>,
                ]}
            >
                <CreateUser
                    onSuccess={() => {
                        setIsModalOpen(false);
                        fetchUsers(currentPage, pageSize, searchValue);
                    }}
                />
            </Modal>
        </>
    );
};

export default UsersList;