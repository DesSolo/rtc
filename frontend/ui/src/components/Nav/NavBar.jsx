import { Flex, Layout, Menu, Typography, Avatar } from "antd";
import {
    ProjectOutlined,
    AuditOutlined,
    UserOutlined
} from '@ant-design/icons';
import { useNavigate } from "react-router-dom";
import {getUsername, hasRole} from "../../utils/storage.js";

const { Text } = Typography;
const { Sider } = Layout;

const NavBar = ({ collapsed }) => {
    const navigate = useNavigate();
    const username = getUsername();

    return (
        <Sider trigger={null} collapsible collapsed={collapsed} style={{ position: 'relative' }}>
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
                    ...(hasRole('admin') ? [{
                        key: '/users',
                        icon: <UserOutlined />,
                        label: 'Users',
                    }]: [])
                ]}
                onClick={(e) => {
                    navigate(e.key);
                }}
            />
            <div
                style={{
                    position: 'fixed',
                    bottom: 0,
                    left: 0,
                    width: collapsed ? 80 : 200,
                    padding: '16px',
                    backgroundColor: '#001529',
                    borderTop: '1px solid #444',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 12,
                    color: 'white',
                    transition: 'width 0.2s',
                    zIndex: 1000
                }}
            >
                <Avatar icon={<UserOutlined />} size="small" />
                {!collapsed && <Text style={{ color: 'white' }}>{username || 'Guest'}</Text>}
            </div>
        </Sider>
    );
};

export default NavBar;