package maquinarefrescs is

    protected type Monitor_Productor_Consumidor is

        -- Inicialitza el monitor amb el nombre de clients
        procedure iniciar (n_clients : Integer);

        -- Entry per al reposador
        entry Reposar (id : Integer; fi_reposador : out boolean);

        -- Entry per als clients que volen agafar un refresc
        entry Agafar
          (n_refresc_agafats      : Integer;
           n_refrescs_a_consumir : Integer;
           nom_client             : String);

        -- Entry per indicar que un client ha acabat
        entry fi_client (clients_restants : out Integer);

    private
        -- Nombre maxim de refrescs que pot contenir la maquina
        max_refrescs_maquina : integer;

        -- Nombre actual de refrescs a la maquina
        n_refrescs_maquina : integer := 0;

        -- Nombre total de clients
        total_clients : integer := 0;

        -- Nombre de clients que han acabat
        clients_acabats : integer := 0;

    end Monitor_Productor_Consumidor;

end maquinarefrescs;
