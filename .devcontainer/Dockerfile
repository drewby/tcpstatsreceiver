FROM mcr.microsoft.com/devcontainers/go:0-1.20-bullseye

# RUN wget https://packages.microsoft.com/config/debian/11/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
# RUN sudo dpkg -i packages-microsoft-prod.deb
# RUN rm packages-microsoft-prod.deb

# RUN sudo apt-get update && \
#   sudo apt-get install -y dotnet-sdk-6.0

ARG USERNAME=vscode
RUN usermod -s /bin/zsh $USERNAME
USER $USERNAME

RUN go install go.opentelemetry.io/collector/cmd/builder@latest
RUN go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@latest

# set up oh-my-zsh
RUN sh -c "$(wget -O- https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" "" --unattended &> /dev/null
ENV ENV="/home/$USERNAME/.ashrc" \
    ZSH=/home/$USERNAME/.oh-my-zsh \
    EDITOR=vi \
    LANG=en_US.UTF-8
RUN printf 'ZSH_THEME="robbyrussell"\nENABLE_CORRECTION="false"\nplugins=(git copyfile extract colorize dotenv encode64 golang)\nsource $ZSH/oh-my-zsh.sh\nexport PATH=$PATH:~/.local/bin' > "/home/$USERNAME/.zshrc"
RUN echo "exec `which zsh`" > "/home/$USERNAME/.ashrc"