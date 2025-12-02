import React from 'react';
import { useSearchParams } from 'react-router-dom';
import { TabPane, Tab } from 'semantic-ui-react';

const Monitor = () => {
  const [searchParams] = useSearchParams();

  const baseUrl = searchParams.get('base_url');

  const protocol = new URL(baseUrl).protocol; // protocol = "http:"
  const hostname = new URL(baseUrl).hostname; // hostname = "127.0.0.1"
  const basePath = `${protocol}//${hostname}:13000`;

  const renderIframe = (title, path) => (
    <TabPane attached={false}>
      <iframe
        style={{ width: '100%', height: '100%' }}
        title={title}
        src={`${basePath}${path}`}
      ></iframe>
    </TabPane>
  );

  const panes = [
    {
      menuItem: 'GPU信息',
      render: () => (
        <TabPane attached={false}>
          {renderIframe(
            'GPU信息',
            '/d/Oxed_c6Wz/nvidia-dcgm-exporter-dashboard?orgId=1&from=now-15m&to=now&timezone=browser&var-instance=dcgm-exporter:9400&var-gpu=$__all&theme=light&kiosk=true'
          )}
        </TabPane>
      ),
    },
    {
      menuItem: '主机详情',
      render: () => (
        <TabPane attached={false}>
          {renderIframe(
            '主机详情',
            '/d/9CWBzd1f0bik001/linuxe4b8bb-e69cba-e8afa6-e68385?orgId=1&from=now-24h&to=now&timezone=browser&var-project=$__all&var-job=linux&var-node=125.69.160.50:19100&var-hostname=2d53376afc16&var-device=eth0&var-maxmount=%2Fmnt%2Fdisk1&var-show_hostname=2d53376afc16&theme=light&kiosk=true'
          )}
        </TabPane>
      ),
    },
  ];

  return (
    <div className='monitor-content'>
      <h2 className='monitor-title'>{searchParams.get('name')}</h2>
      <Tab className='monitor-tab' menu={{ pointing: true }} panes={panes} />
    </div>
  );
};

export default Monitor;
