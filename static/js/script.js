document.addEventListener("DOMContentLoaded", fetchInstances);

const avatarPlaceholder = "/static/img/avatar.png"; // Avatar padrão

let allInstances = []; // Armazena todas as instâncias para pesquisa

/** Busca e exibe as instâncias */
async function fetchInstances() {
    const response = await fetch('/fetchInstances');
    const instances = await response.json();
    allInstances = instances; // Salva os dados para pesquisa
    renderInstances(instances);
}

/** Renderiza as instâncias na tela */
/** Renderiza as instâncias na tela */
function renderInstances(instances) {
    const container = document.getElementById('instancesContainer');
    container.innerHTML = '';

    if (instances.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-warning text-center p-3 d-flex align-items-center justify-content-center">
                    <i class="fas fa-exclamation-circle text-warning me-2" style="font-size: 1.2rem;"></i>
                    <span class="fw-bold" style="font-size: 1rem;color: #000;">Nenhuma instância encontrada!</span>
                </div>
            </div>
        `;
        return;
    }

    const isDarkMode = document.body.classList.contains("dark-mode");
    instances.forEach(async (instance, index) => {
        const profilePic = await validateImage(instance.profilePicUrl);
        
        // Se a imagem não for válida, exibe o ícone do Font Awesome
        let avatarHTML = "";
        if (profilePic === avatarPlaceholder) {
            avatarHTML = `
              <div class="d-flex align-items-center justify-content-center rounded-circle border border-1 frosted-glass-avatar" 
                   style="width:70px; height:70px;">
                <i class="fas fa-user avatar" style="font-size: 2rem;"></i>
              </div>
            `;
          } else {
            avatarHTML = `
              <img src="${profilePic}" alt="Avatar" class="rounded-circle border border-1" 
                   style="width:70px; height:70px; object-fit: cover; background-color: #005c4b;">
            `;
          }
          
        
        const card = document.createElement('div');
        card.classList.add('col-lg-4', 'col-md-6', 'mb-4'); // Responsividade aprimorada
        card.setAttribute('data-name', instance.name.toLowerCase());
        card.setAttribute('data-status', instance.connectionStatus.toLowerCase());
        
        card.innerHTML = `
          <div class="card h-100 rounded-3 shadow-sm custom-card" style="border: 2px solid rgba(1, 110, 90, 0.5); border-top: none;">
            <!-- Área Superior: Avatar centralizado com margin-top negativa -->
            <div class="d-flex justify-content-center" style="margin-top: -35px;">
              ${avatarHTML}
            </div>
            
            <!-- Conteúdo inferior: toda a área abaixo do avatar -->
            <div >
              <!-- Corpo do card com informações -->
              <div class="card-body text-center p-0">
                <h5 class="card-title fw-bold mb-1">${instance.name}</h5>
                <p class="small  mb-1">${instance?.profileName ?? ''}</p>
                <p class="small  mb-3">${instance.ownerJid}</p>
                
                <!-- Grupo de Input para o Token -->
                <div class="input-group my-3 mx-auto" style="max-width: 300px;">
                  <input type="password" id="token-${index}" 
                         value="${instance?.Auth?.token ?? instance?.token ?? ''}" 
                         class="form-control" readonly>
                  <button class="btn" style="background-color: #005c4b; border-color: #005c4b; color: white;" 
                          onclick="toggleToken('token-${index}', 'eye-icon-${index}')" type="button">
                    <i id="eye-icon-${index}" class="fas fa-eye"></i>
                  </button>
                  <button class="btn" style="background-color: #005c4b; border-color: #005c4b; color: white;" 
                          onclick="copyToClipboard('token-${index}')" type="button">
                    <i class="fas fa-copy"></i>
                  </button>
                </div>
              </div>
              
              <!-- Rodapé do card: Status e Botão de Deleção -->
              <div style="border: 1px solid #ced4da;padding: 20px 15px 10px; border-radius: 5px;" class="d-flex justify-content-between align-items-center px-3 py-2">
                <span class="badge ${getStatusBadge(instance.connectionStatus)}">${instance.connectionStatus}</span>
                
               <div class="d-flex">
                   <button onclick="handleConnect('${instance.name}')" class="btn btn-outline-info btn-sm me-2">
                        <i class="fas fa-link"></i> Conect
                    </button>
                  <!-- Botão de Delete -->
                  <button onclick="handleDelete('${instance.name}', this)" class="btn btn-outline-danger btn-sm">
                    <i class="fas fa-trash-alt"></i> Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        `;
        
        container.appendChild(card);
      });
      
    
}


/** Filtra as instâncias pelo nome e status */
function filterInstances() {
    const searchQuery = document.getElementById('search').value.toLowerCase();
    const statusFilter = document.getElementById('statusFilter').value.toLowerCase();

    const filteredInstances = allInstances.filter(instance => {
        const matchesName = instance.name.toLowerCase().includes(searchQuery);
        const matchesStatus = statusFilter === "" || instance.connectionStatus.toLowerCase() === statusFilter;
        return matchesName && matchesStatus;
    });

    renderInstances(filteredInstances);
}

/** Adiciona eventos para filtro e pesquisa */
document.getElementById('search').addEventListener('input', filterInstances);
document.getElementById('statusFilter').addEventListener('change', filterInstances);


document.getElementById('create').addEventListener('click', async (e) => {
    e.preventDefault();
        
    const button = document.getElementById('create');
    const instanceName = document.getElementById('instanceName').value;

    if(!instanceName){
        document.getElementById('error-instance').innerHTML = 'O campo Nome da Instância é obriatório';
        return;
    }

    button.disabled = true;
    button.innerHTML = `Criando ...`;

    const response = await fetch('/createInstance', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `instanceName=${encodeURIComponent(instanceName)}`,
    });

    console.log("Resposta recebida:", response.status); // <- Teste aqui também

    const responseData = await response.json().catch(() => null);

    if (response.ok) {
        showAlert('Instância criada com sucesso!', 'success');
        fetchInstances();
        closeModal();
    } else {
        const errorMessage = responseData?.message?.join("\n") || responseData?.message || "Erro ao criar instância.";
        showAlert(errorMessage, 'error');
    }

    button.disabled = false;
    button.innerHTML = "Criar";
});


/** Alterna a visibilidade do token */
function toggleToken(inputId, eyeIconId) {
    const input = document.getElementById(inputId);
    const eyeIcon = document.getElementById(eyeIconId);

    if (input.type === "password") {
        input.type = "text";
        eyeIcon.classList.replace("fa-eye", "fa-eye-slash");
    } else {
        input.type = "password";
        eyeIcon.classList.replace("fa-eye-slash", "fa-eye");
    }
}

/** Copia o token para a área de transferência */
function copyToClipboard(inputId) {
    const input = document.getElementById(inputId);
    navigator.clipboard.writeText(input.value).then(() => {
        showAlert('Token copiado!', 'success');
    }).catch(() => {
        showAlert('Erro ao copiar!', 'error');
    });
}

/** Valida a imagem e retorna uma alternativa se houver erro */
async function validateImage(url) {
    if (!url) return avatarPlaceholder;

    try {
        const response = await fetch(url, { method: 'HEAD' });
        if (!response.ok || response.status >= 400) {
            return avatarPlaceholder;
        }
        return url;
    } catch {
        return avatarPlaceholder;
    }
}

/** Retorna a classe correta para o status */
function getStatusBadge(status) {
    if (status === "ONLINE") return "status-connected";
    if (status === "OFFLINE") return "status-offline";
    return "badge-connecting";
}

/** Exibe modal de criação */
function openModal() {
    new bootstrap.Modal(document.getElementById('createInstanceModal')).show();
}

function closeModal() {
    const modal = bootstrap.Modal.getInstance(document.getElementById('createInstanceModal'));
    if (modal) {
        modal.hide();
    }
}


/** Confirmação ao deletar */
async function handleDelete(instanceName, button) {
    const confirm = await Swal.fire({
        title: "Tem certeza?",
        text: "Esta ação não pode ser desfeita!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#d33",
        cancelButtonColor: "#3085d6",
        confirmButtonText: "Sim, deletar!",
        cancelButtonText: "Cancelar"
    });

    if (confirm.isConfirmed) {
        button.disabled = true;
        button.innerHTML = `  <div class="mb-2 spinner spinner--danger"><svg viewBox="0 0 50 50"><circle cx="25" cy="25" r="20"></circle></svg></div>`;

        try {
            await deleteInstance(instanceName);
        } catch (error) {
            console.error('Erro ao deletar:', error);
        } finally {
            button.innerHTML = "Delete";
            button.disabled = false;
        }
    }
}

/** Deleta uma instância */
async function deleteInstance(instanceName) {
    try {
        const response = await fetch(`/deleteInstance?instanceName=${instanceName}&force=true`, { method: 'DELETE' });

        // Tenta converter a resposta para JSON
        const responseData = await response.json().catch(() => null);

        console.log("Resposta da API:", response);
        console.log("Dados da resposta:", responseData);

        if (response.ok) {
            showAlert('Instância deletada com sucesso!', 'success');
            fetchInstances();
        } else {
            // Captura a mensagem de erro vinda da API
            const errorMessage = responseData?.message?.join("\n") || responseData?.message || "Erro ao deletar instância.";

            console.error("Erro ao deletar instância:", errorMessage);
            showAlert(errorMessage, 'error');
        }
    } catch (error) {
        console.error("Erro inesperado:", error);
        showAlert("Erro inesperado ao deletar instância.", 'error');
    }
}


/** Exibe alertas bonitos */
function showAlert(message, type) {
    Swal.fire({
        text: message,
        icon: type,
        toast: true,
        position: "top-end",
        showConfirmButton: false,
        timer: 3000
    });
}

/** Alterna entre os modos escuro e claro */
function toggleTheme() {
    const body = document.body;
    const themeIcon = document.getElementById("themeIcon");

    if (body.classList.contains("dark-mode")) {
        body.classList.replace("dark-mode", "light-mode");
        localStorage.setItem("theme", "light-mode");
        themeIcon.classList.replace("fa-moon", "fa-sun"); // Troca ícone para sol
    } else {
        body.classList.replace("light-mode", "dark-mode");
        localStorage.setItem("theme", "dark-mode");
        themeIcon.classList.replace("fa-sun", "fa-moon"); // Troca ícone para lua
    }

    applyThemeToElements();
}

/** Aplica o tema correto aos elementos dinâmicos */
function applyThemeToElements() {
   
}


/** Mantém o tema ao recarregar a página */
document.addEventListener("DOMContentLoaded", () => {
    const savedTheme = localStorage.getItem("theme") || "dark-mode";
    document.body.classList.remove("dark-mode", "light-mode"); // Remove ambas as classes de tema
    document.body.classList.add(savedTheme); // Adiciona a classe de tema salva

    // Ajusta o ícone corretamente ao carregar
    const themeIcon = document.getElementById("themeIcon");
    if (savedTheme === "light-mode") {
        themeIcon.classList.replace("fa-moon", "fa-sun");
    } else {
        themeIcon.classList.replace("fa-sun", "fa-moon");
    }

    // Adiciona evento ao botão de alternância de tema
    document.getElementById("themeToggle").addEventListener("click", toggleTheme);
});

let qrRefreshInterval; // Variável global para armazenar o intervalo

async function fetchQRCode(instanceName) {
  const qrContainer = document.getElementById('qrCodeContainer');
  const errorContainer = document.getElementById('connectError');
  const modalEl = document.getElementById('connectInstanceModal');

  try {
    const response = await fetch(`/instance/connect?instanceName=${encodeURIComponent(instanceName)}`);
    
    if (!response.ok) {
      const errorData = await response.json();
      console.error("Erro recebido:", errorData);
      if (errorData.message && errorData.message.includes("Instance already connected")) {
        const modalInstance = bootstrap.Modal.getInstance(modalEl);
        if (modalInstance) {
          modalInstance.hide();
        }
        fetchInstances();
        return;
      } else {
        throw new Error('Erro ao obter QR Code');
      }
    }
    
    const data = await response.json();
    console.log("QR Code atualizado:", data);
    
    if (data.base64) {
      qrContainer.innerHTML = `<img src="${data.base64}" alt="QR Code para Conectar" class="img-fluid" />`;
    } else {
      qrContainer.innerHTML = '<p>QR Code não disponível.</p>';
    }
    errorContainer.classList.add('d-none');
    errorContainer.textContent = '';
  } catch (error) {
    console.error(error);
    errorContainer.textContent = 'Erro ao obter QR Code.';
    errorContainer.classList.remove('d-none');
    qrContainer.innerHTML = '';
  }
}

async function handleConnect(instanceName) {
  const modalEl = document.getElementById('connectInstanceModal');
  const modal = new bootstrap.Modal(modalEl);
  const qrContainer = document.getElementById('qrCodeContainer');
  const errorContainer = document.getElementById('connectError');
  
  qrContainer.innerHTML = '<p>Carregando QR Code...</p>';
  errorContainer.classList.add('d-none');
  errorContainer.textContent = '';
  
  modal.show();
  
  await fetchQRCode(instanceName);
  
  qrRefreshInterval = setInterval(() => {
    fetchQRCode(instanceName);
  }, 30000);
  
  modalEl.addEventListener('hidden.bs.modal', () => {
    clearInterval(qrRefreshInterval);
  }, { once: true });
}

