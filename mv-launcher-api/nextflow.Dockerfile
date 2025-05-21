FROM ubuntu:22.04

ENV MINICONDA_URL="https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh"
ENV INSTALL_PATH="/nextflow_awscli"
ENV PATH="${INSTALL_PATH}/bin:$PATH"

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    wget \
    jq \
    openjdk-17-jdk \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Install Miniconda and AWS CLI
RUN wget -q $MINICONDA_URL -O miniconda.sh && \
    bash miniconda.sh -b -f -p $INSTALL_PATH && \
    rm miniconda.sh && \
    ${INSTALL_PATH}/bin/conda install -c conda-forge -y awscli && \
    ${INSTALL_PATH}/bin/conda clean -a -y

# Install Nextflow
RUN curl -s https://get.nextflow.io | bash && \
    mv nextflow /usr/local/bin/ && \
    chmod +x /usr/local/bin/nextflow

# Add launcher script
COPY run.sh /run.sh
RUN chmod +x /run.sh

# Set working directory
WORKDIR /workspace

ENTRYPOINT ["/run.sh"]