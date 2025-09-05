import {Flex, Layout, Menu, Typography} from "antd";
import {
  ProjectOutlined,
  AuditOutlined
} from '@ant-design/icons';
import {useNavigate} from "react-router-dom";

const {Text} = Typography;

const { Sider } = Layout;

const NavBar = ({collapsed}) => {
    const navigate = useNavigate();

    return (
        <Sider trigger={null} collapsible collapsed={collapsed}>
            <Flex justify="center" align="center">
                <Text style={{fontSize:34, color: "white", cursor: "pointer"}} onClick={()=>navigate('/')}>RTC</Text>
            </Flex>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['1']}
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
              navigate(e.key)
          }}
        />
      </Sider>
    )
};

export default NavBar;