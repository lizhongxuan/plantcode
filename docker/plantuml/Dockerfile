# 基于官方PlantUML镜像，添加中文字体支持
FROM plantuml/plantuml-server:jetty

# 切换到root用户安装字体
USER root

# 更新包管理器并安装中文字体
RUN apt-get update && \
    apt-get install -y \
    fonts-wqy-zenhei \
    fonts-wqy-microhei \
    fonts-arphic-ukai \
    fonts-arphic-uming \
    fontconfig && \
    fc-cache -fv && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 创建字体配置文件，确保中文字体被正确识别
RUN mkdir -p /usr/share/fonts/truetype/chinese && \
    echo '<?xml version="1.0"?>' > /etc/fonts/local.conf && \
    echo '<!DOCTYPE fontconfig SYSTEM "fonts.dtd">' >> /etc/fonts/local.conf && \
    echo '<fontconfig>' >> /etc/fonts/local.conf && \
    echo '  <alias>' >> /etc/fonts/local.conf && \
    echo '    <family>sans-serif</family>' >> /etc/fonts/local.conf && \
    echo '    <prefer>' >> /etc/fonts/local.conf && \
    echo '      <family>WenQuanYi Zen Hei</family>' >> /etc/fonts/local.conf && \
    echo '      <family>WenQuanYi Micro Hei</family>' >> /etc/fonts/local.conf && \
    echo '    </prefer>' >> /etc/fonts/local.conf && \
    echo '  </alias>' >> /etc/fonts/local.conf && \
    echo '  <alias>' >> /etc/fonts/local.conf && \
    echo '    <family>serif</family>' >> /etc/fonts/local.conf && \
    echo '    <prefer>' >> /etc/fonts/local.conf && \
    echo '      <family>AR PL UKai CN</family>' >> /etc/fonts/local.conf && \
    echo '      <family>AR PL UMing CN</family>' >> /etc/fonts/local.conf && \
    echo '    </prefer>' >> /etc/fonts/local.conf && \
    echo '  </alias>' >> /etc/fonts/local.conf && \
    echo '</fontconfig>' >> /etc/fonts/local.conf

# 设置Java系统属性，确保PlantUML使用正确的字体
ENV JAVA_OPTS="-Dfile.encoding=UTF-8 -Dsun.jnu.encoding=UTF-8 -Djava.awt.headless=true"

# 切换回非root用户（如果原镜像有的话）
USER 1001

# 暴露端口
EXPOSE 8080

# 启动命令保持不变
CMD ["java", "-jar", "/app/plantuml.jar", "-picoweb"]