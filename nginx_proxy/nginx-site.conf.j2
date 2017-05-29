server {
    listen      80;

    location /api/ {
      proxy_set_header Host $host:$server_port;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

      {% if VIDEO_API_SERVICE_HOST is defined %}
      proxy_pass http://{{ VIDEO_API_SERVICE_HOST }}:3000;
      {% else %}
      proxy_pass http://video_api:3000/;
      {% endif %}
    }

    location / {
      proxy_set_header Host $host:$server_port;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

      {% if WEB_CLIENT_SERVICE_HOST is defined %}
      proxy_pass http://{{ WEB_CLIENT_SERVICE_HOST }}:3000;
      {% else %}
      proxy_pass http://web_client:3000/;
      {% endif %}
    }
}
