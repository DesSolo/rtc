import { Flex, Layout, Menu, Typography, Avatar } from "antd";
import {
    ProjectOutlined,
    AuditOutlined,
    UserOutlined // 👈 Добавляем иконку пользователя
} from '@ant-design/icons';
import { useNavigate } from "react-router-dom";
import { getUsername } from "../../utils/storage.js";

const { Text } = Typography;
const { Sider } = Layout;

const NavBar = ({ collapsed }) => {
    const navigate = useNavigate();
    const username = getUsername(); // Получаем имя пользователя

    return (
        <Sider trigger={null} collapsible collapsed={collapsed}>
            {/* Логотип */}
            <Flex justify="center" align="center" style={{ height: 64 }}>
                <Text
                    style={{
                        fontSize: 34,
                        color: "white",
                        cursor: "pointer",
                        fontWeight: "bold"
                    }}
                    onClick={() => navigate('/')}
                >
                    RTC
                </Text>
            </Flex>

            {/* Меню */}
            <Menu
                theme="dark"
                mode="inline"
                defaultSelectedKeys={['/']}
                items={[
                    {
                        key: '/',
                        icon: <ProjectOutlined />,
                        label: 'Projects',
                    },
                    {
                        key: '/audit',
                        icon: <AuditOutlined />,
                        label: 'Audit logs',
                    },
                ]}
                onClick={(e) => {
                    navigate(e.key);
                }}
            />

            {/* Информация о пользователе внизу */}
            <div
                style={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    padding: '16px',
                    backgroundColor: '#001529', // Цвет фона Sider в темной теме
                    borderTop: '1px solid #444',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 12,
                    color: 'white',
                }}
            >
                <Avatar icon={<UserOutlined />} size="small" />
                {!collapsed && <Text style={{ color: 'white' }}>{username || 'Guest'}</Text>}
            </div>
        </Sider>
    );
};

export default NavBar;