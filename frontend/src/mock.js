import Mock from 'mockjs';

Mock.mock('/snatch', {
    'code': 0,
    'msg': 'success',
    'data': {
        'cur_count': 1,
        'max_count': 4,
        "envelope_id": 123,
    }
});

// Mock.mock('/snatch', {
//     'code': 1,
//     'msg': 'failed, some message',
//     'data': {
//         'cur_count': 1,
//         'max_count': 4
//     }
// });

Mock.mock('/open', {
    "code": 0,   // 成功则code=0，否则为其他（请自行定义各类错误）
    "msg": "success",
    "data": {
        "value": 50   // 红包金额，以“分”为单位
    }
});

Mock.mock('/get_wallet_list', {
    'code': 0,
    'msg': 'success',
    'data': {
        "amount": 112,    // 钱包总额，“分”为单位
        "envelope_list": [
            {
                "envelope_id": 123,
                "value": 50,     // 红包面值
                "opened": true,   // 是否已拆开
                "snatch_time": 1634551711     // 红包获取时间，UNIX时间戳
            },
            {
                "envelope_id": 124,
                "opened": false,   // 未拆开的红包不显示value
                "snatch_time": 1634551711
            }
        ]
    }
});
