pipeline {
    agent any
    options {
        retry(1)
    }
    environment {
        project_name = 'goin-web' // 项目名称 每个项目单独设置 需要和docker-compose.yml一致
        credentials_id = '5bd64c05-cabd-44cc-b81e-a76e8af7cfb8'
        project_tar_name = "${project_name}-${BUILD_ID}.tar.gz"
        LEADER_JENKINS = 'maxiongfei' // 审批人jenkins账号
    }
    parameters {
        string(name: 'git_url', defaultValue: 'ssh://git@xxx.git', description: '代码库地址')
        gitParameter(
            branch: '',
            branchFilter: '*',
            defaultValue: 'master',
            description: '发布的分支或tag',
            name: 'branch',
            quickFilterEnabled: false,
            selectedValue: 'NONE',
            sortMode: 'NONE',
            tagFilter: '*',
            type: 'PT_BRANCH_TAG')
        choice(
            name: 'run_env',
            choices: 'test\nprod',
            description: '运行环境：test：测试环境\n prod：生产环境'
        )
        string(name: 'container_base_dir_path', defaultValue: '/data/www', description: '容器内的www目录')
        string(name: 'ssh_user', defaultValue: 'root', description: '待发布目标服务器的ssh用户名')
        string(name: 'ssh_port', defaultValue: '22', description: '待发布目标服务器的ssh端口')
        string(name: 'origin_base_dir_path', defaultValue: '/data/wwwroot', description: '待发布目标服务器的www目录')
        string(name: 'wechat_robot_notice_key', defaultValue: '', description: '企业微信机器人key')
        // 暂时用不上
        booleanParam(name: 'run_migrate', defaultValue: false, description: '是否进行数据库迁移')
        booleanParam(name: 'run_db_seed', defaultValue: false, description: '是否进行数据库数据填充')
    }

    stages {
        // 拉取代码
        stage('拉取代码') {
            steps {
                echo "开始拉取代码  当前发布分支： ${params.branch}  运行环境：${params.run_env}"
                echo "工作目录：${WORKSPACE}"
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: "${params.branch}"]],
                    doGenerateSubmoduleConfigurations: false,
                    extensions: [],
                    gitTool: 'Default',
                    submoduleCfg: [],
                    userRemoteConfigs: [[
                        credentialsId: "${credentials_id}",
                        url: "${params.git_url}"
                    ]]
                ])
                sh "ls -alh"
            }
        }
        stage('打包代码') {
            steps {
                script{
                    echo "开始打包代码"
                    echo "制作docker-compose的env文件"
                    def docker_compose_env = """
                    BUILD_ID=${BUILD_ID}
                    CONTAINER_WWW_PATH=${params.container_base_dir_path}
                    CRONTAB_ON=false
                    COMPOSE_PROJECT_NAME=${project_name}
                    RUN_ENV=${params.run_env}
                    """
                    writeFile file: './deploy/.env', text: docker_compose_env, encoding: 'UTF-8'

                    sh """
                        echo "删除当前目录tar.gz文件"
                        echo "打包文件并查看"
                        rm -rf *.tar.gz
                        tar -zcvf ${project_tar_name} ./* .[!.git]*
                        ls -alh
                    """
                }
            }
        }
        stage('上传代码') {
            steps{
                script{
                     env.Master_Confirm = "No"
                     if (params.run_env == 'prod') {
                         // 通知询问构建
                         def jobInfo=env.JOB_NAME.split("/")
                         def str = "### jenkins发布申请通知\\n >项目名称:`${jobInfo[1]}`\\n >build-id:`${BUILD_ID}` \\n >发布模块--env:`${jobInfo[2]}`--`${params.run_env}` \\n >请开发组长用个人账户登录Jenkins Pipeline页面，同意「`${params.branch}`」分支部署\\n >处理页面:[点击跳转](${env.BUILD_URL}/input)\\n"
                         sh """
                             curl -s -H "Content-Type:application/json" -X POST -d '{"msgtype":"markdown","markdown":{"content":"${str}"}}' https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=${params.wechat_robot_notice_key}
                             curl -s -H "Content-Type:application/json" -X POST -d '{"msgtype":"text","text":{"content":"有新的发布申请，请及时处理;build-id:${BUILD_ID}","mentioned_mobile_list":["18710839146" , "@all"]}}' https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=${params.wechat_robot_notice_key}
                         """

                         // 询问是否同意
                         try {
                             timeout(time: 5, unit: 'MINUTES') {
                                 env.Master_Confirm = input(
                                     message: "是否开始部署-${jobInfo[1]}-${jobInfo[2]}-${params.run_env}-环境代码？",
                                     ok: "确定",
                                     submitter: "${env.LEADER_JENKINS}",
                                     parameters: [
                                         choice(choices: "Yes\nNo\n", description: '开发组长确认是否部署，不同意请选择No!', name: 'Master_Confirm')
                                     ],
                                 )
                             }
                         } catch (err){
                             echo err.toString()
                             env.Master_Confirm='No'
                             error "管理员拒绝了本次发布!"
                         }
                    }else if (params.run_env == 'dev') {
                        echo "发布${params.run_env}环境代码,无需审核"
                        env.Master_Confirm = "Yes"
                    } else if (params.run_env == 'test') {
                        echo "发布${params.run_env}环境代码,无需审核"
                        env.Master_Confirm = "Yes"
                    } else {
                        error "错误的env环境，发布结束"
                        false
                    }
                    if ("${env.Master_Confirm}".contains("Yes")) {
                        echo "部署服务器中。。。"
                        def blue = deploy_code()
                        parallel blue
                        echo "部署服务器完成。。。"
                    } else {
                        false
                    }
                }
            }
        }
    }
	post {
        success{
            script {
                def jobInfo=env.JOB_NAME.split("/")
                def str = "### 代码发布完成通知「成功」\\n >项目名称:${jobInfo[1]}\\n >build-id:${BUILD_ID} \\n  >发布结果:`成功`\\n >发布分支:${params.branch}\\n >查看部署详情:[点击跳转](${env.BUILD_URL}/console)\\n"
                sh """
                    curl -H "Content-Type:application/json" -X POST -d '{"msgtype":"markdown","markdown":{"content":"${str}"}}' https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=${params.wechat_robot_notice_key}
                """
            }
        }

        failure{
            script {
                // 所有分支出错都需要通知
                def jobInfo=env.JOB_NAME.split("/")
                def str = "### 代码发布完成通知「失败」\\n >项目名称:${jobInfo[1]}\\n >build-id:${BUILD_ID} \\n >发布结果:`失败`\\n >发布分支:${params.branch}\\n >查看部署详情:[点击跳转](${env.BUILD_URL}/console)\\n"
                sh """
                    curl -H "Content-Type:application/json" -X POST -d '{"msgtype":"markdown","markdown":{"content":"${str}"}}' https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=${params.wechat_robot_notice_key}
                """
            }
        }

        unstable{
            script {
                // 所有分支不稳定都需要通知
                script {
                    // 所有分支出错都需要通知
                    def jobInfo=env.JOB_NAME.split("/")
                    def str = "### 代码发布完成通知「失败」\\n >项目名称:${jobInfo[1]}\\n >build-id:${BUILD_ID}\\n  >发布结果:`失败`\\n >发布分支:${params.branch}\\n >查看部署详情:[点击跳转](${env.BUILD_URL}/console)\\n"
                    sh """
                        curl -H "Content-Type:application/json" -X POST -d '{"msgtype":"markdown","markdown":{"content":"${str}"}}' https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=${params.wechat_robot_notice_key}
                    """
                }
            }
        }
	}
}

def deploy_code() {
    script {
        def parallelDeploy  = [:]

        def hosts = [
            "test": [
                "server1",
            ],
            "prod": [
                "server1",
                "server2",
            ]
        ]

        def domains = [
            "test": "goin-web.com",
            "prod": "goin-web.com"
        ]
        def consul_names = [
            "test": "goin-web.com",
            "prod": "goin-web.com"
        ]

        def current_hosts = hosts.get(run_env)
        def current_domain = domains.get(run_env)
        def current_consul_name = consul_names.get(run_env)

        def origin_www_path="${params.origin_base_dir_path}/${current_domain}" // 项目服务器的www目录
        def docker_volume_path="${origin_www_path}/volume" // 挂载目录，一般用于日志目录挂载
        def jenkins_build_history_path="${origin_www_path}/jenkins-build/" // jenkins发布历史记录
        def docker_deploy_path="${origin_www_path}/deploy/${BUILD_ID}" // 目标服务器的docker部署目录
        def deploy_info="${current_domain}-${BUILD_ID}" // 项目名

        for (int i = 0; i < current_hosts.size(); i++) {
            def ip = current_hosts[i]
            def j = i
            def ssh = "ssh -p ${params.ssh_port} ${params.ssh_user}@${ip}"
            parallelDeploy["deploy-task-${ip}-${i}"] = {
                sh """
                    ${ssh} "mkdir -p -m 777 ${origin_www_path} ${jenkins_build_history_path} ${docker_deploy_path} ${docker_volume_path}"
                    scp -P ${ssh_port} ${project_tar_name} ${ssh_user}@${ip}:${jenkins_build_history_path}
                    ${ssh} "cd ${jenkins_build_history_path} && tar xf ${project_tar_name} -C ${docker_deploy_path}"
                    ${ssh} "cd ${docker_deploy_path} && echo VOLUME_PATH=${docker_volume_path} >> ./deploy/.env"
                    ${ssh} "cd ${docker_deploy_path} && echo DOCKER_DEPLOY_PATH=${docker_deploy_path} >> ./deploy/.env"
                    ${ssh} "cd ${docker_deploy_path} && echo CONSUL_SERVICE_NAME=${current_consul_name} >> ${docker_deploy_path}/deploy/.env"

                    if [ $j -eq 0 ]
                    then
                        echo ${ip}开启计划任务
                        ${ssh} "sed -i 's/CRONTAB_ON=false/CRONTAB_ON=true/g' ${docker_deploy_path}/deploy/.env"
                    fi

                    echo ${ip}销毁容器并启动新的容器
                    ${ssh} "cd ${docker_deploy_path}/deploy && docker-compose -f docker-compose.yml down && docker-compose -f docker-compose.yml up -d --build"

                    if [ "${run_migrate}" = "true" ] && [ $j -eq 0 ]
                    then
                        echo 执行数据库迁移
                        ssh ${ssh_user}@${ip} -p ${ssh_port} "cd ${docker_deploy_path}/deploy && docker-compose exec -T --user=www-data ${project_name} php artisan migrate";
                    fi
                    if [ "${run_db_seed}" = "true" ] && [ $j -eq 0 ]
                    then
                        echo 执行数据库迁移
                        ssh ${ssh_user}@${ip} -p ${ssh_port} "cd ${docker_deploy_path}/deploy && docker-compose exec -T --user=www-data ${project_name} php artisan db:seed";
                    fi
                """
            }
        }
        return parallelDeploy
    }
}
