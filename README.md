Seguir os seguintes passos para execução do projeto:

1. Gerar a imagem por meio do comando make docker-build na raiz do projeto
2. Executar o comando para realizar os testes de stress:
    docker run stress-test --url=https://urlsite.com --concurrency=X --requests=Y
3. No final da execução, um relatório será gerado e exibido no terminal